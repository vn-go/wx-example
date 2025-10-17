// Đảm bảo đã thêm using cho các namespace cần thiết
using hrm.Data;
using hrm.Models;
using Microsoft.EntityFrameworkCore;
using Pomelo.EntityFrameworkCore.MySql.Infrastructure;
using System.Collections;
var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddOpenApi();

// --- Bắt đầu: Cấu hình MySQL và DbContext ---
// 1. Lấy chuỗi kết nối từ appsettings.json
var connectionString = builder.Configuration.GetConnectionString("DefaultConnection");

// 2. Đăng ký DbContext với MySQL
builder.Services.AddDbContext<HrmDbContext>(options =>
    options.UseMySql(
        connectionString,
        ServerVersion.AutoDetect(connectionString),
        mySqlOptions =>
        {
            // XÓA BỎ DÒNG GÂY LỖI:
            // .CharSetBehavior(CharSetBehavior.NeverAppend) 
        }
    )
);
// --- Kết thúc: Cấu hình MySQL và DbContext ---

var app = builder.Build();
// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.MapOpenApi();
}

app.UseHttpsRedirection();



// --- Endpoint mới: Lấy Danh sách Role ---
app.MapGet("/api/roles", async (HrmDbContext db) =>
{
    try
    {
        // 1. Kiểm tra cache
        if (Cache.Roles != null)
        {
            return Results.Ok(Cache.Roles);
        }

        // 2. Truy vấn DB và Caching
        var rolesFromDb = await db.SysRoles.ToListAsync();
        Cache.Roles = rolesFromDb; // Gán vào biến static trong lớp Cache

        return Results.Ok(rolesFromDb);
    }
    catch (Exception ex)
    {
        // Xử lý lỗi (ví dụ: lỗi kết nối DB)
        // Trong môi trường production, bạn nên log lỗi chi tiết hơn
        return Results.Problem(
            title: "Lỗi truy vấn cơ sở dữ liệu",
            detail: ex.Message,
            statusCode: StatusCodes.Status500InternalServerError);
    }
})
.WithName("GetListOfRoles")
.WithOpenApi(); // Thêm vào Swagger/OpenAPI
app.Run();


