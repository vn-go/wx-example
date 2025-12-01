export function findViewPath(ele: HTMLElement): string | undefined {
    let currentElement: HTMLElement | null = ele;

    // Lặp chừng nào currentElement còn tồn tại và chưa phải là body
    while (currentElement && currentElement.tagName !== 'BODY') {
        const viewPath = currentElement.getAttribute("view-path");

        // 1. Kiểm tra thuộc tính trên phần tử hiện tại
        if (viewPath) {
            return viewPath;
        }

        // 2. Chuyển lên phần tử cha kế tiếp
        currentElement = currentElement.parentElement;
    }

    // Nếu không tìm thấy 'view-path' trên cây DOM cho đến thẻ <body>
    return undefined;
}