var builder = WebApplication.CreateBuilder(args);

// Bật Swagger cực nhẹ, không cần package ngoài (từ .NET 8+ đã có sẵn)
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

var app = builder.Build();

// Bật Swagger UI (chạy ngon trên .NET 10 preview)
app.UseSwagger();
app.UseSwaggerUI();

app.MapGet("/", () => Results.Json(new { message = "Hello from .NET 10!", time = DateTime.UtcNow }))
   .WithName("Root")
   .WithOpenApi();

app.MapPost("/users", (User user) => Results.Created($"/users/{user.Id}", user))
   .WithName("CreateUser")
   .WithOpenApi();

app.Run();

record User(int Id, string Name, int Age, string[] Skills);
//dotnet new console -f net10.0 --no-restore