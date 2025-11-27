// File: PersonJsonContext.cs
using System.Text.Json.Serialization;

[JsonSerializable(typeof(Person))]
[JsonSerializable(typeof(Address))]
[JsonSerializable(typeof(int[]))]
[JsonSerializable(typeof(string[]))]
internal partial class PersonJsonContext : JsonSerializerContext
{
}