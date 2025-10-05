using Microsoft.EntityFrameworkCore;
using hrm.Models;

namespace hrm.Data
{
    public class HrmDbContext : DbContext
    {
        // 1. Tạo DbSet cho bảng sys_roles
        public DbSet<SysRole> SysRoles { get; set; } = default!;

        public HrmDbContext(DbContextOptions<HrmDbContext> options)
            : base(options)
        {
        }

        // 2. Tùy chỉnh ánh xạ nếu cần (ví dụ: ràng buộc UNIQUE KEY)
        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            // Cấu hình các ràng buộc UNIQUE KEY trong code-first (Tùy chọn)
            modelBuilder.Entity<SysRole>()
                .HasIndex(r => r.Code)
                .IsUnique();

            modelBuilder.Entity<SysRole>()
                .HasIndex(r => r.Name)
                .IsUnique();

            modelBuilder.Entity<SysRole>()
                .HasIndex(r => r.RoleId)
                .IsUnique();
        }
    }
}