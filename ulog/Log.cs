using UnityEngine;
using System.Collections;
using System;
using System.Collections.Generic;
using System.Text;

namespace ulog
{
    public enum LogLevel
    {
        Info = 0,
        Warning,
        Error,
        NoLog,
    }

    public static class Log
    {
        private static LogLevel LogLevel = LogLevel.Info;

        public static void LogInfo(string msg, params object[] args)
        {
            if (LogLevel <= LogLevel.Info)
            {
                if (args.Length == 0)
                    Debug.Log(msg);
                else
                {
                    Debug.Log(string.Format(msg, args));
                }
            }
        }
        public static void LogTODO(string msg, params object[] args)
        {
            if (LogLevel <= LogLevel.Info)
            {
                msg = string.Format("TODO:{0}", msg);
                if (args.Length == 0)
                    Debug.Log(msg);
                else
                    Debug.Log(string.Format(msg, args));
            }
        }
        public static void LogWarning(string msg, params object[] args)
        {
            if (LogLevel <= LogLevel.Warning)
            {
                if (args.Length == 0)
                    Debug.LogWarning(msg);
                else
                    Debug.LogWarning(string.Format(msg, args));
            }
        }
        public static void LogError(string msg, params object[] args)
        {
            if (LogLevel <= LogLevel.Error)
            {
                if (args.Length == 0)
                    Debug.LogError(msg);
                else
                    Debug.LogError(string.Format(msg, args));
            }
        }
        public static void LogException(Exception ex)
        {
            if (LogLevel <= LogLevel.Error)
            {
                Debug.LogException(ex);
            }
        }

        #region log time
        private static float time = 0;
        public static void LogDeltaTime(string label)
        {
            if (time == 0)
            {
                time = Time.time;
            }
            LogInfo(label + ":" + (Time.time - time));
            time = Time.time;
        }
        #endregion
    }
}
