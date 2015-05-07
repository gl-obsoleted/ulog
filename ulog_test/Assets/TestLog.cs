using UnityEngine;
using System.Collections;
using ulog;
using System;

public class TestLog : MonoBehaviour {

	// Use this for initialization
	void Start () {
        Log.LogInfo("test log info: {0} {1} {2}", 0, 3.5f, "foo");
        Log.LogTODO("test TODO.");
        Log.LogWarning("test warning.");
        Log.LogError("test error.");
        Log.LogAssert(false, "test assert");
        Log.LogException(new Exception());
	}
	
	// Update is called once per frame
	void Update () {
	
	}
}
