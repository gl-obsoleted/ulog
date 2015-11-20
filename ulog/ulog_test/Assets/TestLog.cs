using UnityEngine;
using System.Collections;
using System;

public class TestLog : MonoBehaviour {

    public bool LoggingIntoFile = true;
    public bool PrintTestLogs = true;
    public int DaysKeeping = -1;  // would cleanup if older than 'n' days

    private LogService _logServ;

	// Use this for initialization
	void Start () 
    {
        LogService.UserDefinedMemBufSize = 8 * 1024; // setup a 8KB memory buffer

        int days = DaysKeeping == 0 ? 1 : DaysKeeping;
        _logServ = new LogService(LoggingIntoFile, days, true);

        if (PrintTestLogs)
        {
            Log.Info(null);
            Log.Info("test log info: {0} {1} {2}", 0, 3.5f, "foo");

            Log.TODO(null);
            Log.TODO("test TODO.");
            Log.TODO("test TODO param {0}, {1}.", 1, "abc");

            Log.Trace(null);
            Log.Trace("test Trace.");
            Log.Trace("test Trace param {0}, {1}.", 1, "abc");

            for (int i = 0; i < 20; i++)
                Log.Info("repeat AAA.");

            Log.InfoEx(null, null);
            Log.InfoEx(null, this);
            Log.InfoEx("test context object.", null);
            Log.InfoEx("test context object.", this);

            for (int i = 0; i < 20; i++)
                Log.Info("repeat BBB.");

            Log.Warning("test warning.");
            Log.Error("test error.");
            Log.Assert(false, "test assert");

            Log.Exception(new Exception("foo"));
            Log.Exception(new Exception("bar"));

            for (int i = 0; i < 30; i++)
                Log.Error(new Exception("Oops! Error."));

            for (int i = 0; i < 30; i++)
                Log.Exception(new Exception("Oops! Exception."));

            foreach (var item in LogUtil.InMemoryExceptions)
            {
                Log.Info("[LogUtil.InMemoryExceptions]: {0}", item);
            }

            foreach (var item in LogUtil.InMemoryErrors)
            {
                Log.Info("[LogUtil.InMemoryErrors]: {0}", item);
            }

            StartCoroutine(RunTimedLogging());
        }
	}

    IEnumerator RunTimedLogging()
    {
        Log.DeltaTime("init");
        yield return new WaitForSeconds(0.2f);
        Log.DeltaTime("sleep for 0.2s");
        yield return new WaitForSeconds(0.3f);
        Log.DeltaTime("sleep for 0.3s");
        yield return new WaitForSeconds(0.4f);
        Log.DeltaTime("sleep for 0.4s");
        yield return new WaitForSeconds(0.3f);

        for (int i = 0; i < 100; i++)
            Log.Info("repeat info.");
    }

    // Update is called once per frame
	void Update () {
	
	}

    void OnDestroy()
    {
        if (_logServ != null)
        {
            _logServ.Dispose();
        }
    }
}
