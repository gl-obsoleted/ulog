using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.IO;
using System.Linq;
using System.Text;
using UnityEngine;

public class LogBuffer
{
    public const int KB = 1024;
    public const int BufSize = 16 * KB;

    public byte[] Buf = new byte[BufSize];
    public int BufWrittenBytes = 0;

    public bool Receive(string content)
    {
        byte[] bytes = Encoding.Default.GetBytes(content);
        if (BufWrittenBytes + bytes.Length > BufSize)
            return false;

        Buffer.BlockCopy(bytes, 0, Buf, BufWrittenBytes, bytes.Length);
        BufWrittenBytes += bytes.Length;
        return true;
    }

    public void Clear()
    {
        BufWrittenBytes = 0;
    }
}


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

    public LogService(bool logIntoFile, int oldLogsKeptDays) // '-1' means keeping all logs without any erasing
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

        if (oldLogsKeptDays > 0)
        {
            try
            {
                CleanupLogsOlderThan(oldLogsKeptDays);
            }
            catch (Exception e)
            {
                Log.Exception(e);
                Log.Error("Cleaning up logs ({0}) failed.", oldLogsKeptDays);
            }
        }
    }

    public void Dispose()
    {
        FlushLogWriting();

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

    public void WriteLog(string content, LogType type)
    {
        // write directly if larger than buffer
        if (Encoding.Default.GetByteCount(content) > LogBuffer.BufSize)
        {
            if (_logWriter != null)
            {
                _logWriter.Write(content);
            }
        }

        // write into buffer 
        if (type == LogType.Error || !_memBuf.Receive(content))
        {
            // flush into file when buffer is full
            FlushLogWriting();

            _memBuf.Receive(content);
        }
    }

    public void FlushLogWriting()
    {
        if (_logWriter != null)
        {
            _logWriter.Write(Encoding.Default.GetString(_memBuf.Buf, 0, _memBuf.BufWrittenBytes));
        }
        _memBuf.Clear();
    }

    private void CleanupLogsOlderThan(int days)
    {
        DateTime timePointForDeleting = DateTime.Now.Subtract(TimeSpan.FromDays(days));
        string timeStrForDeleting = SysUtil.FormatDateAsFileNameString(timePointForDeleting);

        DirectoryInfo logDirInfo = new DirectoryInfo(SysUtil.CombinePaths(Application.persistentDataPath, "log"));
        DirectoryInfo[] dirsByDate = logDirInfo.GetDirectories();
        List<string> toBeDeleted = new List<string>();
        foreach (var item in dirsByDate)
        {
            //Log.Info("[COMPARING]: {0}, {1}", item.Name, timeStrForDeleting);
            if (string.CompareOrdinal(item.Name, timeStrForDeleting) <= 0)
            {
                toBeDeleted.Add(item.FullName);
                //Log.Info("[TO_BE_DELETED]: {0}", item.FullName);
            }
        }

        foreach (var item in toBeDeleted)
        {
            Directory.Delete(item, true);
            Log.Info("[ Log Cleanup ]: {0}", item);
        }
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
            WriteLog(string.Format("{0:0.00} {1}: {2}\r\n", Time.realtimeSinceStartup, type, condition), type);

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
            Log.Exception(ex); // this should at least print to Unity Editor (but may skip the file writing due to earlier writing failure)            	
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

    private LogBuffer _memBuf = new LogBuffer();

    public static string LastLogFile { get; set; }
}
