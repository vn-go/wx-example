// src/lib/utils/dom.ts
/**
 * Tìm element đầu tiên có thuộc tính slot="<slotName>"
 * Đợi liên tục cho đến khi tìm thấy hoặc hết timeout.
 *
 * @param slotName - Tên slot (VD: "body", "header")
 * @param timeout - Thời gian chờ tối đa (ms), mặc định 5000ms
 * @returns Promise<HTMLElement | null>
 *
 * @example
 * const body = await findBySlotName('body', 10000);
 * if (body) body.style.background = 'yellow';
 */
export async function findBySlotName(
    ele: HTMLElement | undefined,
    slotName: string,
    timeout: number = 5000
): Promise<HTMLElement | undefined> {
    const startTime = Date.now();

    return new Promise((resolve) => {
        if (!ele) resolve(undefined)
        const interval = setInterval(() => {
            // Tìm element có thuộc tính slot đúng tên
            const element = ele?.querySelector<HTMLElement>(`[slot="${slotName}"]`);

            if (element) {
                clearInterval(interval);
                resolve(element);
            } else if (Date.now() - startTime >= timeout) {
                clearInterval(interval);
                console.warn(`[findBySlotName] Timeout: Không tìm thấy slot="${slotName}" sau ${timeout}ms`);
                resolve(undefined);
            }
        }, 100); // Kiểm tra mỗi 100ms

        // Kiểm tra ngay lập tức lần đầu (tránh delay 100ms)
        const immediate = document.querySelector<HTMLElement>(`[slot="${slotName}"]`);
        if (immediate) {
            clearInterval(interval);
            resolve(immediate);
        }
    });
}