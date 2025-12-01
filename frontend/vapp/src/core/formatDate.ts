/**
 * formatDate
 * @param dateValue: Date | string | number
 * @param format: string, vd: "dd/MM/yyyy", "yyyy-MM-dd HH:mm:ss", "MM-dd-yyyy"
 */
export function formatDate(dateValue: Date | string | number, format = "dd/MM/yyyy"): string {
    if (!dateValue) return "";

    const date = dateValue instanceof Date ? dateValue : new Date(dateValue);
    if (isNaN(date.getTime())) return "";

    const map: Record<string, string> = {
        "dd": String(date.getDate()).padStart(2, "0"),
        "d": String(date.getDate()),
        "MM": String(date.getMonth() + 1).padStart(2, "0"),
        "M": String(date.getMonth() + 1),
        "yyyy": String(date.getFullYear()),
        "yy": String(date.getFullYear()).slice(-2),
        "HH": String(date.getHours()).padStart(2, "0"),
        "H": String(date.getHours()),
        "hh": String(date.getHours() % 12 || 12).padStart(2, "0"),
        "h": String(date.getHours() % 12 || 12),
        "mm": String(date.getMinutes()).padStart(2, "0"),
        "m": String(date.getMinutes()),
        "ss": String(date.getSeconds()).padStart(2, "0"),
        "s": String(date.getSeconds()),
        "a": date.getHours() < 12 ? "AM" : "PM"
    };

    // Replace tất cả token trong format
    return format.replace(/yyyy|yy|MM|M|dd|d|HH|H|hh|h|mm|m|ss|s|a/g, matched => map[matched]);
}
