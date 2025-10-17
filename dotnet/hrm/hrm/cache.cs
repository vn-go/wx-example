using hrm.Data;
using hrm.Models;
using Microsoft.EntityFrameworkCore;
using Pomelo.EntityFrameworkCore.MySql.Infrastructure;
using System.Collections;
public static class Cache
{
    // Lỗi CS0106 được khắc phục vì biến static nằm trong một lớp static
    public static List<SysRole>? Roles = null;
}