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
        int days = DaysKeeping == 0 ? 1 : DaysKeeping;
        _logServ = new LogService(LoggingIntoFile, days);

        if (PrintTestLogs)
        {
            Log.Info("test log info: {0} {1} {2}", 0, 3.5f, "foo");

            for (int i = 0; i < 30; i++)
                Log.Info("repeat AAA.");

            Log.TODO("test TODO.");

            for (int i = 0; i < 50; i++)
                Log.Info("repeat BBB.");

            Log.Warning("test warning.");
            Log.Error("test error.");
            Log.Assert(false, "test assert");
            Log.Exception(new Exception());

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
