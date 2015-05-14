using UnityEngine;
using System.Collections;
using GameCommon;
using System;

public class TestLog : MonoBehaviour {

    public bool LoggingIntoFile = true;
    public bool PrintTestLogs = true;

    private LogService _logServ;

	// Use this for initialization
	void Start () 
    {
        _logServ = new LogService(LoggingIntoFile);

        if (PrintTestLogs)
        {
            Log.Info("test log info: {0} {1} {2}", 0, 3.5f, "foo");
            Log.TODO("test TODO.");
            Log.Warning("test warning.");
            Log.Error("test error.");
            Log.Assert(false, "test assert");
            Log.Exception(new Exception());
        }
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
