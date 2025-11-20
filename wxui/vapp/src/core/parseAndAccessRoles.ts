import { WindowWiget } from "./draggableWindowWiget";

/**
 * Chuyển đổi chuỗi HTML thành DocumentFragment và truy cập các phần tử con theo thuộc tính role.
 * * @param htmlContent Chuỗi HTML cần phân tích cú pháp.
 * @returns Một đối tượng chứa các phần tử header và body (hoặc null nếu không tìm thấy).
 */
function parseAndAccessRoles(htmlContent: string): WindowWiget {
    // 1. Khởi tạo DOMParser
    const parser = new DOMParser();

    // 2. Phân tích cú pháp chuỗi HTML
    // documentFragment sẽ chứa các nodes được tạo ra từ chuỗi HTML
    const doc = parser.parseFromString(htmlContent, 'text/html');

    // Quan trọng: Vì chuỗi HTML của bạn chỉ là một fragment (không có <html>, <body>), 
    // chúng ta cần lấy phần tử đầu tiên trong body của document được parse.
    // Lấy phần tử container chính:
    const modalContainer = doc.body.firstChild as HTMLElement;

    if (!modalContainer) {
        console.error("Không tìm thấy container chính.");
        return null;
    }

    // 3. Truy cập các phần tử theo thuộc tính 'role'
    // Sử dụng querySelector cho thuộc tính attribute [role='value']
    const headerElement = modalContainer.querySelector('[role="header"]') as HTMLElement | null;
    const bodyElement = modalContainer.querySelector('[role="body"]') as HTMLElement | null;
    const titleElement = headerElement.querySelector('[role="title"]') as HTMLElement
    const btnCloseElement = headerElement.querySelector('[role="close"]') as HTMLElement //role="close"

    // 4. Trả về kết quả
    return {
        container: modalContainer,
        header: headerElement,
        body: bodyElement,
        title: titleElement,
        closeBtn: btnCloseElement
    };
}

export default parseAndAccessRoles