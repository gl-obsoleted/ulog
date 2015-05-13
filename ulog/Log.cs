﻿using UnityEngine;
using System;

namespace GameCommon
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
        public static LogLevel LogLevel = LogLevel.Info;

        public static void Info(string msg, params object[] args)
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
        public static void TODO(string msg, params object[] args)
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
        public static void Warning(string msg, params object[] args)
        {
            if (LogLevel <= LogLevel.Warning)
            {
                if (args.Length == 0)
                    Debug.LogWarning(msg);
                else
                    Debug.LogWarning(string.Format(msg, args));
            }
        }
        public static void Error(string msg, params object[] args)
        {
            if (LogLevel <= LogLevel.Error)
            {
                if (args.Length == 0)
                    Debug.LogError(msg);
                else
                    Debug.LogError(string.Format(msg, args));
            }
        }
        public static void Exception(Exception ex)
        {
            if (LogLevel <= LogLevel.Error)
            {
                Debug.LogException(ex);
            }
        }

        public static void Assert(bool condition)
        {
            if (LogLevel <= LogLevel.Error)
            {
                Assert(condition, string.Empty, true);
            }
        }

        public static void Assert(bool condition, string assertString)
        {
            if (LogLevel <= LogLevel.Error)
            {
                Assert(condition, assertString, false);
            }
        }

        public static void Assert(bool condition, string assertString, bool pauseOnFail)
        {
            if (!condition && LogLevel <= LogLevel.Error)
            {
                Debug.LogError("assert failed! " + assertString);

                if (pauseOnFail)
                    Debug.Break();
            }
        }

        #region log time
        private static float time = 0;
        public static void DeltaTime(string label)
        {
            if (time == 0)
            {
                time = Time.time;
            }
            Info(label + ":" + (Time.time - time));
            time = Time.time;
        }
        #endregion
    }
}