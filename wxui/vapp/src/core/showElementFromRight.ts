/**
 * Animate an HTMLElement to "show" with a smooth entrance:
 * 1. Start from the right edge of the viewport, vertically centered
 * 2. Then slide to the exact center of the viewport (both horizontally and vertically)
 * 
 * @param element The HTML element to animate
 * @param duration Total duration of the entire animation in milliseconds (default: 800ms)
 * @param easing Timing function (default: ease-out)
 */
function showElementFromRight(element: HTMLElement, onFinish: () => {}, duration: number = 300, easing: string = 'ease-out'): Promise<undefined> {
    const ret = new Promise<undefined>((resolve, reject) => {
        // Ensure the element is positioned absolutely or fixed for transform to work properly
        element.style.position = 'fixed';
        element.style.top = '50%';        // Will be adjusted with translateY later
        element.style.left = '100%';      // Start completely off-screen to the right
        element.style.transform = 'translateY(-50%)'; // Vertically center it
        element.style.opacity = '0';      // Start invisible
        element.style.zIndex = '9999';    // Make sure it's on top (adjust if needed)

        // Force reflow to ensure initial state is applied before animation
        void element.offsetWidth;

        // Phase 1: First, move from right edge → still on the right but visible + vertically centered
        // Phase 2: Then move to the true center of the screen
        const keyframes = [
            {
                // Starting point: off-screen right, vertically centered, invisible
                left: '100%',
                transform: 'translateY(-50%)',
                opacity: '0',
            },
            {
                // Mid point: still on the right edge but now visible (quick fade-in + small slide)
                left: 'calc(100% - 100px)', // Slightly enter the viewport (adjust as you like)
                transform: 'translateY(-50%)',
                opacity: '1',
                offset: 0.4, // 40% of total duration
            },
            {
                // Final point: exact center of viewport
                left: '50%',
                transform: 'translateX(-50%) translateY(-50%)',
                opacity: '1',
            }
        ];

        const options: KeyframeAnimationOptions = {
            duration,
            easing,
            fill: 'forwards', // Keep the final state after animation ends
        };

        // Run the animation
        const animation = element.animate(keyframes, options);

        animation.onfinish = () => {
            // Commit the final styles so future changes won't break
            element.style.left = '50%';
            element.style.top = '50%';
            element.style.transform = 'translateX(-50%) translateY(-50%)';
            element.style.opacity = '1';
            onFinish();
            resolve(undefined);
        };
    });
    // Optional: Resolve when animation finishes
    return ret;
}
export default showElementFromRight