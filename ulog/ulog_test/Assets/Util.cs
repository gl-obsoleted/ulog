using UnityEngine;
using System.Collections;
using System.IO;
using System.Security.Cryptography;
using System.Text;

public delegate void onFileProcessed(string filepath);

public class Util
{
    // Process all files in the directory passed in, recurse on any directories  
    // that are found, and process the files they contain. 
    public static void ProcessDirectory(string targetDirectory, onFileProcessed handleFile)
    {
        // Process the list of files found in the directory. 
        string[] fileEntries = Directory.GetFiles(targetDirectory);
        foreach (string fileName in fileEntries)
            handleFile(fileName);

        // Recurse into subdirectories of this directory. 
        string[] subdirectoryEntries = Directory.GetDirectories(targetDirectory);
        foreach (string subdirectory in subdirectoryEntries)
            ProcessDirectory(subdirectory, handleFile);
    }

    public static long GetFileSize(string filepath)
    {
        try
        {
            return new FileInfo(filepath).Length;
        }
        catch (System.Exception)
        {
            return 0;            
        }
    }

    public static string GetFileMD5(string filepath)
    {
        try
        {
            StringBuilder b = new StringBuilder();
            using (var md5 = MD5.Create())
            {
                using (var stream = File.OpenRead(filepath))
                {
                    byte[] digest = md5.ComputeHash(stream);
                    for (int i = 0; i < digest.Length; i++)
                    {
                        b.AppendFormat("{0:x}", digest[i]);
                    }
                }
            }
            return b.ToString();
        }
        catch (System.Exception)
        {

            return "";
        }
    }
}
