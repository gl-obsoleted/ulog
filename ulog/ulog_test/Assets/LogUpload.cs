using UnityEngine;
using System.Collections;
using System.Collections.Generic;
using System.Text;
using System;
using System.IO;

public class LogUpload : MonoBehaviour
{
    public bool EnableLogUploading = true;
    public string LogServerAddr = "localhost";
    public int LogServerPort = 13080;

	private string _ticket = "";
    private List<string> _uploadPreparing = new List<string>();
    private List<string> _uploaded = new List<string>();

	void Start ()
    {
        StartCoroutine(PerformUploadProcess());
    }

    public IEnumerator PerformUploadProcess()
    {
        // reset vars
        _ticket = "";
        _uploadPreparing.Clear();
        _uploaded.Clear();

        Log.Info("Querying ticket...");
        yield return StartCoroutine(QueryTicket());
        if (_ticket.Length == 0)
        {
            Log.Info("Querying ticket failed.");
            yield break;
        }
        Log.Info("Ticket : {0}.", _ticket);

        Log.Info("Verify files...");
        yield return StartCoroutine(VerifyLogs());
        if (_uploadPreparing.Count == 0)
        {
            Log.Info("No file needs uploading.");
            yield break;
        }
        Log.Info("Verifying files Done, {0} files queued.", _uploadPreparing.Count);

        Log.Info("Uploading files...");
        yield return StartCoroutine(UploadLogs());
        Log.Info("Uploading files Done, {0} files uploaded.", _uploaded.Count);
    }

    public IEnumerator QueryTicket()
	{
        byte[] bytesToEncode = Encoding.UTF8.GetBytes("1|2|3|4");
        string encodedText = Convert.ToBase64String(bytesToEncode);

        WWWForm form = new WWWForm();
        form.AddField("user_info", encodedText);
        WWW w = new WWW(BuildURL("query_ticket"), form);
        yield return w;
        if (!string.IsNullOrEmpty(w.error))
        {
            Log.Info(w.error);
        }
        else
        {
            _ticket = System.Text.Encoding.Default.GetString(w.bytes); 
        }
	}

    public IEnumerator VerifyLogs()
	{
        List<string> logfiles = new List<string>();

        string logDir = LogUtil.CombinePaths(Application.persistentDataPath, "log");
        Util.ProcessDirectory(logDir, (path) => { logfiles.Add(path); });

        Dictionary<string, string> loginfo = new Dictionary<string, string>();
        foreach (var item in logfiles)
        {
            long size = Util.GetFileSize(item);
            string md5 = Util.GetFileMD5(item);
            if (size != 0 && md5.Length != 0)
            {
                loginfo.Add(item, string.Format("{0}|{1}", size, md5));
            }
        }

        WWWForm form = new WWWForm();
        foreach (var item in loginfo)
            form.AddField(item.Key, item.Value);
        WWW w = new WWW(BuildURL("validate_files", _ticket), form);
        yield return w;
        if (!string.IsNullOrEmpty(w.error))
        {
            Log.Info(w.error);
        }
        else
        {
            string[] neededFiles = System.Text.Encoding.Default.GetString(w.bytes).Split('|');
            _uploadPreparing = new List<string>(neededFiles);
        }
    }

    public IEnumerator UploadLogs()
	{
        foreach (var item in _uploadPreparing)
        {
            Log.Info("Uploading {0}.", item);

            WWWForm form = new WWWForm();
            form.AddBinaryData("file_list",
                System.IO.File.ReadAllBytes(item),
                new FileInfo(item).Name,
                "multipart/form-data");
            WWW w = new WWW(BuildURL("upload_resource", _ticket), form);
            yield return w;
            if (!string.IsNullOrEmpty(w.error))
            {
                Log.Info(w.error);
            }
            else
            {
                Log.Info("Uploaded {0} ({1}).", item, System.Text.Encoding.Default.GetString(w.bytes));
                _uploaded.Add(item);
            }
        }
        yield break;
    }

    private string BuildURL(string verb, string ticket = "")
    {
        if (ticket.Length == 0)
	    {
            return string.Format("http://{0}:{1}/{2}", LogServerAddr, LogServerPort, verb);
	    }
        else
	    {
            return string.Format("http://{0}:{1}/{2}?ticket={3}", LogServerAddr, LogServerPort, verb, ticket);
	    }
    }
}
