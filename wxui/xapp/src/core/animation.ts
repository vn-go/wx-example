// src/lib/utils/animation.ts
/**
 * Slide animation: formEle (right → left), toEle (right → center)
 * Uses Tailwind classes + CSS transition.
 *
 * @param formEle - Current element (will slide out to left)
 * @param toEle   - Next element (will slide in from right)
 * @param duration - Animation duration in ms (default: 400)
 */
export function slideFromRightToLeft(
    newEle: HTMLElement | undefined,
    container: HTMLElement | undefined,
    duration: number = 400
): Promise<void> {
    return new Promise((resolve) => {
        if (!newEle || !container) {
            resolve();
            return;
        }
        // 1. Lấy old child
        const oldChild = container.children[0] as HTMLElement;
        if (!oldChild || oldChild === newEle) {
            resolve();
            return;
        }

        // 2. Đảm bảo container có position relative
        const originalPosition = container.style.position;
        container.style.position = 'relative';
        container.style.overflow = 'hidden';

        // 3. Đặt cả 2 element absolute
        oldChild.style.position = 'absolute';
        oldChild.style.top = '0';
        oldChild.style.left = '0';
        oldChild.style.width = '100%';
        oldChild.style.height = '100%';
        oldChild.style.zIndex = '10';
        oldChild.style.transition = `all ${duration}ms ease-in-out`;

        newEle.style.position = 'absolute';
        newEle.style.top = '0';
        newEle.style.right = '0';
        newEle.style.width = '100%';
        newEle.style.height = '100%';
        newEle.style.zIndex = '20';
        newEle.style.transform = 'translateX(100%)';
        newEle.style.transition = `all ${duration}ms ease-in-out`;

        // 4. Force reflow
        void container.offsetWidth;

        // 5. Bắt đầu animation
        oldChild.style.transform = 'translateX(-100%)';
        newEle.style.transform = 'translateX(0)';

        // 6. Kết thúc
        const onEnd = () => {
            // Xóa old child
            oldChild.remove();

            // Reset newEle
            newEle.style.position = '';
            newEle.style.top = '';
            newEle.style.right = '';
            newEle.style.width = '';
            newEle.style.height = '';
            newEle.style.transform = '';
            newEle.style.zIndex = '';
            newEle.style.transition = '';

            // Khôi phục container
            container.style.position = originalPosition;
            container.style.overflow = '';

            resolve();
        };

        setTimeout(onEnd, duration);
    });
}