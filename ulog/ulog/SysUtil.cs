using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;

public class SysUtil
{
    public static string CombinePaths(params string[] paths)
    {
        if (paths == null)
        {
            throw new ArgumentNullException("paths");
        }

        return paths.Aggregate(Path.Combine);
    }

    public static string FormatDateAsFileNameString(DateTime dt)
    {
        return string.Format("{0:0000}-{1:00}-{2:00}", dt.Year, dt.Month, dt.Day);
    }

    public static string FormatTimeAsFileNameString(DateTime dt)
    {
        return string.Format("{0:00}-{1:00}-{2:00}", dt.Hour, dt.Minute, dt.Second);
    }
}
