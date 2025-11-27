using System.Text.Json;
using System.Linq;
using System.Text.Json.Serialization;

Console.WriteLine("Đang tạo JSON lớn ~1.9MB...");

// Tạo đúng mảng 100_000 phần tử
var bigArray = Enumerable.Range(0, 100_000).ToArray();

var person = new Person
{
    Name = "Nguyễn Văn A",
    Age = 28,
    IsDeveloper = true,
    Skills = ["C#", "Go", "Kubernetes", ".NET 10"],
    Address = new Address { City = "Hà Nội", District = "Cầu Giấy" },
    BigArray = bigArray
};

// DÙNG SOURCE GEN – đây là chìa khóa!
var json = JsonSerializer.Serialize(person, PersonJsonContext.Default.Person);

Console.WriteLine($"JSON size: {json.Length / 1024.0:F1} KB");

// Warmup
_ = JsonSerializer.Deserialize(json, PersonJsonContext.Default.Person);

// Đo tốc độ 100 lần
var sw = System.Diagnostics.Stopwatch.StartNew();
const int loops = 100;
for (int i = 0; i < loops; i++)
{
    _ = JsonSerializer.Deserialize(json, PersonJsonContext.Default.Person);
}
sw.Stop();

Console.WriteLine($".NET 10 parse {loops} lần: {sw.ElapsedMilliseconds} ms");
Console.WriteLine($"Trung bình mỗi lần: {sw.ElapsedMilliseconds / (double)loops:F3} ms");

var p = JsonSerializer.Deserialize(json, PersonJsonContext.Default.Person);
Console.WriteLine($"Name: {p?.Name}, BigArray count: {p?.BigArray.Length}");

[JsonSerializable(typeof(Person))]
[JsonSerializable(typeof(Address))]
[JsonSerializable(typeof(int[]))]
[JsonSerializable(typeof(string[]))]
internal partial class PersonJsonContext : JsonSerializerContext
{
}

public class Person
{
    public string Name { get; set; } = "";
    public int Age { get; set; }
    public bool IsDeveloper { get; set; }
    public string[] Skills { get; set; } = [];
    public Address Address { get; set; } = new();
    public int[] BigArray { get; set; } = [];
}

public class Address
{
    public string City { get; set; } = "";
    public string District { get; set; } = "";
}