using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.IO;
using System.Linq;
using System.Text;
using UnityEngine;

public class LogEventArgs : EventArgs
{
    public LogEventArgs(int seqID, LogType type, string content, float time)
    {
        SeqID = seqID;
        LogType = type;
        Content = content;
        Time = time;
    }

    public int SeqID = 0;
    public LogType LogType;
    public string Content = "";
    public float Time = 0.0f;
}

public delegate void LogTargetHandler(object sender, LogEventArgs args);

public class LogService : IDisposable
{
    public event LogTargetHandler LogTargets;

    public LogService(bool logIntoFile)
    {
        RegisterCallback();

        if (logIntoFile)
        {
            try
            {
                DateTime dt = DateTime.Now;

                string logDir = SysUtil.CombinePaths(Application.persistentDataPath, "log", SysUtil.FormatDateAsFileNameString(dt));
                Directory.CreateDirectory(logDir);

                string logPath = Path.Combine(logDir, SysUtil.FormatDateAsFileNameString(dt) + '_' + SysUtil.FormatTimeAsFileNameString(dt) + ".txt");

                _logWriter = new FileInfo(logPath).CreateText();
                _logWriter.AutoFlush = true;
                _logPath = logPath;

                Log.Info("'Log Into File' enabled, file opened successfully. ('{0}')", _logPath);
                LastLogFile = _logPath;
            }
            catch (System.Exception ex)
            {
                Log.Info("'Log Into File' enabled but failed to open file.");
                Log.Exception(ex);
            }
        }
        else
        {
            Log.Info("'Log Into File' disabled.");
            LastLogFile = "";
        }
    }

    public void Dispose()
    {
        Log.Info("destroying log service...");

        if (_logWriter != null)
        {
            _logWriter.Close();
        }

#if UNITY_5_0
        Application.logMessageReceivedThreaded -= OnLogReceived;
#endif

        _disposed = true;
    }

    private void RegisterCallback()
    {
#if UNITY_5_0
        Application.logMessageReceivedThreaded += OnLogReceived;
#else
        Application.RegisterLogCallbackThreaded(OnLogReceived);
#endif
    }

    private void OnLogReceived(string condition, string stackTrace, LogType type)
    {
        if (_disposed || _isWriting)
            return;

        _isWriting = true;

        ++_seqID;

        switch (type)
        {
            case LogType.Assert:
                _assertCount++;
                break;
            case LogType.Error:
                _errorCount++;
                break;
            case LogType.Exception:
                _exceptionCount++;
                break;

            case LogType.Warning:
            case LogType.Log:
            default:
                break;
        }

        try
        {
            if (_logWriter != null)
            {
                _logWriter.WriteLine("{0:0.00} {1}: {2}", Time.realtimeSinceStartup, type, condition);
            }

            foreach (LogTargetHandler Caster in LogTargets.GetInvocationList())
            {
                ISynchronizeInvoke SyncInvoke = Caster.Target as ISynchronizeInvoke;
                LogEventArgs args = new LogEventArgs(_seqID, type, condition, Time.realtimeSinceStartup);

                if (SyncInvoke != null && SyncInvoke.InvokeRequired)
                    SyncInvoke.Invoke(Caster, new object[] { this, args });
                else
                    Caster(this, args);
            }
        }
        catch (System.Exception ex)
        {
            Log.Exception(ex); // this should at least print to Unity Editor (but skip the file writing due to earlier writing failure)            	
        }

        _isWriting = false;
    }

    private string _logPath;
    private StreamWriter _logWriter;

    private ushort _seqID = 0;
    private int _assertCount = 0;
    private int _errorCount = 0;
    private int _exceptionCount = 0;

    private bool _disposed = false;
    private bool _isWriting = false;

    public static string LastLogFile { get; set; }
}
