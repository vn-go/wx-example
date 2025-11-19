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
    async load(componentPath: string, containerId?: string) {
        // 1. Dynamic import Vue component

        const htmlContent = await this.htmlModules[`../${this.htmlLayout}`]();
        const container = document.createElement('div')
        container.innerHTML = htmlContent;
        let childEle = container.children[0];
        document.body.appendChild(childEle)


        const componentLoader = this.componentModules[`../${componentPath}.vue`]
        console.log(componentLoader)
        const Component = (await componentLoader()).default;

        const app = createApp(Component);
        app.mount(childEle)// loi
        /*
        modal.ts:28 [Vue warn]: Component provided template option but runtime compilation is not supported in this build of Vue. Configure your bundler to alias "vue" to "vue/dist/vue.esm-bundler.js". 
  at <App>
        */

        return { container, app }
    }
}
export default Modal