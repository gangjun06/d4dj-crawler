// Refference: https://github.com/GEEKiDoS/D4DJ-Tools/tree/master/AssetTool
using MessagePack;
using MessagePack.Resolvers;
using Newtonsoft.Json;
using Newtonsoft.Json.Converters;
using System;
using System.Collections.Generic;
using System.IO;
using System.Reflection;
using System.Security.Cryptography;

namespace D4DJ_Tool
{
	class Program
	{
		static T DeserializeMsgPack<T>(byte[] decrypted)
		{
			var options = MessagePackSerializerOptions.Standard.WithCompression(MessagePackCompression.Lz4Block);
			return MessagePackSerializer.Deserialize<T>(decrypted, options);
		}

		static void DecryptMaster(FileInfo inputFile, byte[] decrypted)
		{
			var options = MessagePackSerializerOptions.Standard.WithCompression(MessagePackCompression.Lz4Block);

			File.WriteAllText(
			    inputFile.FullName.Replace(".msgpack", ".json"),
			    MessagePackSerializer.ConvertToJson(decrypted, options)
			);
		}


		static string DumpToJson(object obj)
		{
			return JsonConvert.SerializeObject(obj, Formatting.Indented, new StringEnumConverter());
		}


		static void ProcessFileSystemEntry(FileSystemInfo fileSystemInfo)
		{
			if (fileSystemInfo is FileInfo fileInfo)
			{

				var fStream = fileInfo.OpenRead();

				// transform the string into bytes
				byte[] fileData = new byte[fStream.Length];
				// reading the data
				fStream.Read(fileData, 0, fileData.Length);

				if (fileInfo.Name.EndsWith("Master.msgpack"))
				{
					DecryptMaster(fileInfo, fileData);
				}
				else if (fileInfo.Name.EndsWith("ResourceList.msgpack"))
				{
					Console.WriteLine($"Dumping ResourceList...");

					var result = DeserializeMsgPack<Dictionary<string, (int, int)>>(fileData);

					File.WriteAllText(
					    fileInfo.FullName.Replace(".msgpack", ".json"),
					    DumpToJson(result)
					);

				}
			}
		}

		static void Main(string[] args)
		{
			foreach (var arg in args)
			{
				if (File.Exists(arg))
				{
					ProcessFileSystemEntry(new FileInfo(arg));
				}
			}

			Console.WriteLine("Successful!");
		}
	}
}