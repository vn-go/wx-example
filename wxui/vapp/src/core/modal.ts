import { createApp } from 'vue';
import showElementFromRight from './showElementFromRight';
export class ModalInstance {
    private renderEle: HTMLElement;
    private componentPath: string;
    private data?: any;
    private width: number;
    private height: number;
    private templatePath: string;
    private htmlModules: any;
    private componentModules: any;
    private rootEle: HTMLElement;
    constructor(rootEle: HTMLElement, htmlModules, componentModules, templatePath: string, componentPath: string, data?: any) {
        this.htmlModules = htmlModules;
        this.componentModules = componentModules;

        this.templatePath = templatePath;
        this.componentPath = componentPath;
        this.width = 0;
        this.height = 0;
        this.rootEle = rootEle;
        this.data = data;
    }
    setData(data: any): ModalInstance {
        this.data = data;
        return this;
    }
    setSize(width: number, height: number): ModalInstance {
        this.width = width;
        this.height = height;
        return this;
    }
    async render() {
        const htmlContent = await this.htmlModules[`../${this.templatePath}`]();
        const container = document.createElement('div');
        container.innerHTML = htmlContent;
        let childEle = container.children[0];
        childEle.setAttribute("ui-id", this.componentPath);
        this.rootEle.appendChild(childEle);
        showElementFromRight(childEle as HTMLElement);
        const componentLoader = this.componentModules[`../${this.componentPath}.vue`]
        const Component = (await componentLoader()).default;
        const app = createApp(Component, this.data);
        app.mount(childEle)
        childEle.children[0].setAttribute("ui-id", this.componentPath);

        return { childEle, app }
    }
}
class Modal {
    htmlLayout: string;

    htmlModules: any;
    componentModules: any;
    constructor(htmlLayout?: string) {

        this.htmlModules = (import.meta as any).glob('../../src/**/*.html', { as: 'raw' });
        this.htmlLayout = htmlLayout;
        this.componentModules = (import.meta as any).glob('../../src/**/*.vue');



    }
    load(rootEle: HTMLElement, componentPath: string, data?: any): ModalInstance {

        let ret = new ModalInstance(rootEle, this.htmlModules, this.componentModules, this.htmlLayout, componentPath, data);

        return ret;
    }
}
export default Modal
