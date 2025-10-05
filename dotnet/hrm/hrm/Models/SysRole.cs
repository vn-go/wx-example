using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace hrm.Models
{
    [Table("sys_roles")] // Ánh xạ tới tên bảng trong MySQL
    public class SysRole
    {
        // Khóa chính và Identity (AUTO_INCREMENT)
        [Key]
        [Column("id", TypeName = "bigint")]
        public long Id { get; set; }

        // Mặc định EF Core sẽ ánh xạ PascalCase (RoleId) sang snake_case (role_id)
        // nhưng nên dùng [Column] để đảm bảo chính xác.
        [Column("role_id")]
        [MaxLength(36)] // varchar(36)
        public string RoleId { get; set; } = string.Empty;

        [Column("code")]
        [MaxLength(50)] // varchar(50) - UNIQUE KEY
        public string Code { get; set; } = string.Empty;

        [Column("name")]
        [MaxLength(50)] // varchar(50) - UNIQUE KEY
        public string Name { get; set; } = string.Empty;

        [Column("description")]
        [MaxLength(200)] // varchar(200)
        public string Description { get; set; } = string.Empty;

        [Column("created_on")]
        public DateTime CreatedOn { get; set; }

        [Column("modified_on")]
        public DateTime? ModifiedOn { get; set; } // datetime NULLABLE nên dùng DateTime?

        [Column("created_by")]
        [MaxLength(50)] // varchar(50)
        public string CreatedBy { get; set; } = string.Empty;

        [Column("is_active")]
        public bool IsActive { get; set; } // tinyint(1) được ánh xạ thành bool
    }
}