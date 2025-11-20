import { createApp } from 'vue';
import draggableWindowWiget from './draggableWindowWiget';
import parseAndAccessRoles from './parseAndAccessRoles';
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
    setSize(width?: number, height?: number): ModalInstance {
        if (width) this.width = width;
        if (height) this.height = height;
        return this;
    }
    maximize(): ModalInstance {
        this.width = -1;
        this.height = -1;
        return this;
    }
    async render() {
        const htmlContent = await this.htmlModules[`../${this.templatePath}`]();
        const parserDOM = parseAndAccessRoles(htmlContent);
        draggableWindowWiget(parserDOM);
        // const container = document.createElement('div');
        // container.innerHTML = htmlContent;
        let childEle = parserDOM.container as HTMLElement;
        parserDOM.container.setAttribute("ui-id", this.componentPath);
        this._applySize(parserDOM.container)
        this.rootEle.appendChild(parserDOM.container);
        const maxWith = document.body.getBoundingClientRect().width;
        const maxHeight = document.body.getBoundingClientRect().height
        parserDOM.container.style.maxWidth = `${maxWith}px`;
        parserDOM.container.style.maxHeight = `${maxHeight}px`;
        //document.body.appendChild(parserDOM.container);

        let app = undefined;
        await showElementFromRight(childEle as HTMLElement, async () => {
            const componentLoader = this.componentModules[`../${this.componentPath}.vue`]
            const Component = (await componentLoader()).default;
            app = createApp(Component, this.data);
            app.mount(parserDOM.body)
            childEle.children[0].setAttribute("ui-id", this.componentPath);
        });


        return { childEle, app }
    }
    private _applySize(ele: Element): HTMLElement {
        let w = this.width;
        let h = this.height;
        if (this.width == -1) {
            const ret = document.body.getBoundingClientRect();
            w = ret.width;
        }
        if (this.height == -1) {
            const ret = document.body.getBoundingClientRect();
            h = ret.height;
        }
        let htmlEle = ele as HTMLElement
        if (w > 0) {
            htmlEle.style.width = `${w}px`;
        }
        if (h > 0) {
            htmlEle.style.height = `${h}px`;
        }


        return htmlEle;
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
