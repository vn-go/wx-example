// useSlideAnimation.ts

export function useSlideAnimation(targetRef: any) {

    const slideFromRight = (duration = 500) => {
        const el = targetRef.value;
        if (!el) return;

        // set initial state off-screen
        el.style.transition = `transform ${duration}ms ease-out`;
        el.style.transform = 'translateX(100%)';

        // force browser render first frame
        requestAnimationFrame(() => {
            requestAnimationFrame(() => {
                el.style.transform = 'translateX(0)';
            });
        });

        // reset transition after animation ends
        const onEnd = () => {
            el.style.transition = '';
            el.removeEventListener('transitionend', onEnd);
        };
        el.addEventListener('transitionend', onEnd);
    };

    return { slideFromRight };
}
