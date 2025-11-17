// src/lib/utils/dynamic-render.ts
import type { ComponentType } from 'svelte';
import { mount } from 'svelte';

const componentMap = import.meta.glob('./../views/**/*.svelte', { eager: false });

/**
 * Renders a Svelte component and returns the instance.
 *
 * @param virtualPath - e.g., 'auth/login'
 * @param target - DOM element
 * @param props - Props
 * @returns { success: true, instance, destroy } or { success: false, error }
 */
export async function renderDynamicComponent<T = any>(
    virtualPath: string,
    target: HTMLElement | undefined,
    props: Record<string, any> = {}
): Promise<
    | { success: true; instance: T; ins: any, ele: HTMLElement | undefined, destroy: () => void }
    | { success: false; error: string }
> {
    if (!target) {
        return {
            success: false,
            error: 'target is undefined'
        };
    }
    if (virtualPath.startsWith('/')) {
        virtualPath = virtualPath.substring(1, virtualPath.length);
    }
    const relativePath = `../views/${virtualPath}.svelte`;
    const importFn = componentMap[relativePath];

    if (!importFn) {
        return { success: false, error: `Component not found: ${virtualPath}` };
    }

    try {
        target.setAttribute("page-path", virtualPath);
        const module = (await importFn() as any);
        const Component: ComponentType = module.default;

        if (!Component) {
            return { success: false, error: 'No default export' };
        }

        const instance = mount(Component, { target, props }) as T;
        await new Promise(resolve => setTimeout(resolve, 0));

        const root = target.querySelector(`[view-id$="${virtualPath}"]`) || target.firstElementChild;

        let newEle = target.children[target.children.length - 1];
        const instanceOfComponent = (newEle as any)?.__instance as T;
        newEle.setAttribute("view-id", virtualPath);
        return {
            success: true,
            instance,
            ins: instanceOfComponent,
            ele: newEle as any,
            destroy: () => {
                if (typeof (instance as any)?.destroy === 'function') {
                    (instance as any).destroy();
                }
            }
        };
    } catch (error) {
        return {
            success: false,
            error: error instanceof Error ? error.message : 'Render failed'
        };
    }
}