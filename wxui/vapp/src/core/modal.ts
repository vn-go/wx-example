import { createApp } from "vue";
class Modal {
    htmlLayout: string;

    htmlModules: any;
    componentModules: any;
    constructor(htmlLayout?: string) {

        this.htmlModules = (import.meta as any).glob('../../src/**/*.html', { as: 'raw' });
        this.htmlLayout = htmlLayout;
        this.componentModules = (import.meta as any).glob('../../src/**/*.vue');


    }
    async load(componentPath: string, data?: any) {
        // 1. Dynamic import Vue component

        const htmlContent = await this.htmlModules[`../${this.htmlLayout}`]();
        const container = document.createElement('div')
        container.innerHTML = htmlContent;
        let childEle = container.children[0];
        document.body.appendChild(childEle)


        const componentLoader = this.componentModules[`../${componentPath}.vue`]
        console.log(componentLoader)
        const Component = (await componentLoader()).default;

        const app = createApp(Component, data);
        app.mount(childEle)

        return { container, app }
    }
}
export default Modal