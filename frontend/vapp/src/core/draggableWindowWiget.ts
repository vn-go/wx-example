export type WindowWiget = {
    container: HTMLElement, // a whole conatiner
    header: HTMLElement, // indicator for dragging
    body: HTMLElement,
    title: HTMLElement,
    closeBtn: HTMLElement
}
/**
 * Makes the window widget draggable using its header element.
 * * @param widget The WindowWiget object containing the container and header elements.
 */
function draggableWindowWiget(widget: WindowWiget) {
    // 1. Initialize variables for tracking drag state
    let isDragging = false;
    let offsetX: number; // Mouse X coordinate relative to the container's left edge
    let offsetY: number; // Mouse Y coordinate relative to the container's top edge

    // 2. Set initial CSS position style requirement for dragging
    // The container must be positioned absolutely or fixed for 'left' and 'top' to work.
    if (widget.container.style.position !== 'fixed' && widget.container.style.position !== 'absolute') {
        widget.container.style.position = 'absolute';
    }

    // --- MOUSE DOWN (Start Dragging) ---
    const onMouseDown = (e: MouseEvent) => {
        // Only start dragging with the left mouse button (button 0)
        if (e.button !== 0) return;

        // Prevent the default browser drag behavior for elements (e.g., images)
        e.preventDefault();

        isDragging = true;

        // Get the initial position of the mouse relative to the container's top-left corner
        const containerRect = widget.container.getBoundingClientRect();

        // Calculate offset (mouse position - container position)
        offsetX = e.clientX - containerRect.left;
        offsetY = e.clientY - containerRect.top;

        // Optionally, add a class to the container during drag for visual feedback (e.g., cursor change)
        widget.container.classList.add('is-dragging');

        // IMPORTANT: Attach listeners to the whole document to handle mouse movements
        // even if the cursor moves off the header/container area.
        document.addEventListener('mousemove', onMouseMove);
        document.addEventListener('mouseup', onMouseUp);
    };

    // --- MOUSE MOVE (Dragging in Progress) ---
    const onMouseMove = (e: MouseEvent) => {
        if (!isDragging) return;

        // Calculate the new position based on the mouse cursor's current position 
        // minus the initial offset captured in onMouseDown.
        let newX = e.clientX - offsetX;
        let newY = e.clientY - offsetY;

        // Apply the new position to the container style
        widget.container.style.left = `${newX}px`;
        widget.container.style.top = `${newY}px`;

        // Optional: Prevent text selection while dragging
        e.preventDefault();
    };

    // --- MOUSE UP (Stop Dragging) ---
    const onMouseUp = () => {
        if (!isDragging) return;

        isDragging = false;

        // Remove the dragging class
        widget.container.classList.remove('is-dragging');

        // IMPORTANT: Remove event listeners from the document to stop tracking movement
        document.removeEventListener('mousemove', onMouseMove);
        document.removeEventListener('mouseup', onMouseUp);
    };

    // 3. Attach the initial listener to the header element
    widget.header.addEventListener('mousedown', onMouseDown);

    // Optional: Add a style to the header to indicate it's draggable
    widget.header.style.cursor = 'grab';

    // Optional: Add event listener to close button (assuming you want it hooked up here)
    widget.closeBtn.addEventListener('click', () => {
        widget.container.remove();
    });
}
export default draggableWindowWiget