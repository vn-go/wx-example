// src/lib/utils/animation.ts
/**
 * Page transition: old child slides left out, new element slides in from right.
 * Like flipping paper pages.
 *
 * @param newEle     - New element to slide in
 * @param container  - Parent container (has old child as first child)
 * @param duration   - Animation time (ms), default 500
 */
export async function slidePageTransition(
    newEle: HTMLElement | undefined,
    container: HTMLElement | undefined,
    duration: number = 500
): Promise<void> {
    return new Promise((resolve) => {
        if (!newEle || !container) return;
        // 1. Lấy old child (child đầu tiên)
        const oldChild = Array.from(container.children).find(el => el !== newEle) as HTMLElement | undefined;

        //const oldChild = container.firstElementChild as HTMLElement;
        if (!oldChild) {
            container.appendChild(newEle);
            resolve();
            return;
        }

        // 2. Chuẩn bị: đặt cả 2 element vào container
        // Ẩn overflow để không thấy phần trượt ra
        const originalOverflow = container.style.overflow;
        container.style.overflow = 'hidden';
        container.style.position = 'relative';

        // Đặt old child về vị trí ban đầu
        oldChild.style.position = 'absolute';
        oldChild.style.top = '0';
        oldChild.style.left = '0';
        oldChild.style.width = '100%';
        oldChild.style.transition = `transform ${duration}ms ease-in-out, opacity ${duration}ms ease-in-out`;

        // Đặt new element ở bên phải
        newEle.style.position = 'absolute';
        newEle.style.top = '0';
        newEle.style.right = '0';
        newEle.style.width = '100%';
        newEle.style.transform = 'translateX(100%)';
        newEle.style.opacity = '0';
        newEle.style.transition = `transform ${duration}ms ease-in-out, opacity ${duration}ms ease-in-out`;

        // Thêm newEle vào container
        container.appendChild(newEle);

        // Force reflow
        void container.offsetWidth;

        // 3. Bắt đầu animation
        // Old: trượt sang trái + mờ
        oldChild.style.transform = 'translateX(-100%)';
        oldChild.style.opacity = '0';

        // New: trượt từ phải vào giữa
        newEle.style.transform = 'translateX(0)';
        newEle.style.opacity = '1';

        // 4. Kết thúc
        const onEnd = () => {
            // Xóa old child
            // oldChild.remove();
            Array.from(container.children).forEach(el => {
                if (el !== newEle) el.remove();
            });

            // Reset newEle về bình thường
            newEle.style.position = '';
            newEle.style.top = '';
            newEle.style.right = '';
            newEle.style.width = '';
            newEle.style.transform = '';
            newEle.style.opacity = '';
            newEle.style.transition = '';

            // Khôi phục container
            container.style.overflow = originalOverflow;
            container.style.position = '';

            resolve();
        };

        setTimeout(onEnd, duration);
    });
}