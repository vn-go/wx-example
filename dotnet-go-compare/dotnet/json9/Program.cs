// Program.cs
using System.Text.Json;
using System.Linq;

Console.WriteLine("Đang tạo JSON lớn ~575KB (100k phần tử)...");

var bigArray = Enumerable.Range(0, 100_000).ToArray();

var person = new Person
{
    Name = "Nguyễn Văn A",
    Age = 28,
    IsDeveloper = true,
    Skills = ["C#", "Go", ".NET 9"],
    Address = new Address { City = "Hà Nội", District = "Cầu Giấy" },
    BigArray = bigArray
};

// Dùng source-gen → AOT chạy ngon
var json = JsonSerializer.Serialize(person, PersonJsonContext.Default.Person);

Console.WriteLine($"JSON size: {json.Length / 1024.0:F1} KB");

// Warmup
_ = JsonSerializer.Deserialize(json, PersonJsonContext.Default.Person);

var sw = System.Diagnostics.Stopwatch.StartNew();
const int loops = 100;
for (int i = 0; i < loops; i++)
{
    _ = JsonSerializer.Deserialize(json, PersonJsonContext.Default.Person);
}
sw.Stop();

Console.WriteLine($".NET 9 AOT parse {loops} lần: {sw.ElapsedMilliseconds} ms");
Console.WriteLine($"Trung bình mỗi lần: {sw.ElapsedMilliseconds / (double)loops:F3} ms");

var p = JsonSerializer.Deserialize(json, PersonJsonContext.Default.Person); //<--The name 'PersonJsonContext' does not exist in the current contextCS0103

Console.WriteLine($"Name: {p?.Name}, BigArray count: {p?.BigArray.Length}");

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