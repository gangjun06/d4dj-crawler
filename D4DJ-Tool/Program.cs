using MessagePack;
using MessagePack.Resolvers;
using Newtonsoft.Json;
using Newtonsoft.Json.Converters;
using System;
using System.Collections.Generic;
using System.IO;
using System.Reflection;
using System.Security.Cryptography;
using D4DJ_Tools.Masters;

namespace D4DJ_Tools
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
			var typeName = inputFile.Name.Replace(".msgpack", "");
			var targetType = MasterTypes.GetDeserializeType(typeName);

			var options = MessagePackSerializerOptions.Standard.WithCompression(MessagePackCompression.Lz4Block);

			if (targetType == null)
			{
				Console.WriteLine($"Not supported master {typeName}.");
				File.WriteAllText(
				    inputFile.FullName.Replace(".msgpack", ".json"),
				    MessagePackSerializer.ConvertToJson(decrypted, options)
				);
			}
			else
			{
				var result = MessagePackSerializer.Deserialize(targetType, decrypted, options);
				File.WriteAllText(
				    inputFile.FullName.Replace(".msgpack", ".json"),
				    DumpToJson(result)
				);
			}
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
				else if (fileInfo.Name.StartsWith("chart_"))
				{
					try
					{
						object result = null;

						// Check if this is chart common data
						if (fileInfo.Name.EndsWith("0"))
							result = DeserializeMsgPack<ChartCommonData>(fileData);
						else
							result = DeserializeMsgPack<ChartData>(fileData);

						var options = MessagePackSerializerOptions.Standard.WithCompression(MessagePackCompression.Lz4Block);

						File.WriteAllText(
						    fileInfo.FullName + ".json",
						    DumpToJson(result)
						);
					}
					catch (Exception ex)
					{
						Console.WriteLine($"Failed to dump chart: {ex.Message}");
					}
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
		}
	}
}