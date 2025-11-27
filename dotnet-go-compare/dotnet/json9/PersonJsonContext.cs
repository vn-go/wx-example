// PersonJsonContext.cs
using System.Text.Json.Serialization;

[JsonSerializable(typeof(Person))]
[JsonSerializable(typeof(Address))]
internal partial class PersonJsonContext : JsonSerializerContext
{
}