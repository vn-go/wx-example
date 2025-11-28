import PrimeVue from 'primevue/config';
import { App, createApp } from 'vue';
import draggableWindowWiget from './draggableWindowWiget';
import parseAndAccessRoles from './parseAndAccessRoles';
import showElementFromRight from './showElementFromRight';
class UiWatcher {
    app: any;
    container: any;
    observer: MutationObserver;
    constructor(app, container) {
        this.app = app;
        this.container = container;

        this.observer = new MutationObserver(() => {
            if (!document.body.contains(container)) {
                this.app.unmount();
                this.observer.disconnect();
            }
        });

        this.observer.observe(document.body, { childList: true, subtree: true });
    }
}
class rootWatcher {

    observer: MutationObserver;

    constructor(rootEle, childEle) {
        this.observer = new MutationObserver(() => {
            if (!document.body.contains(rootEle)) {
                childEle.remove();
                this.observer.disconnect();
            }
        });

        this.observer.observe(document.body, { childList: true, subtree: true });
    }
}
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
    private _title?: string;
    private ownerViewPath: string;
    constructor(ownerViewPath: string, rootEle: HTMLElement, htmlModules, componentModules, templatePath: string, componentPath: string, data?: any) {
        this.htmlModules = htmlModules;
        this.componentModules = componentModules;
        this.ownerViewPath = ownerViewPath

        this.templatePath = templatePath;
        this.componentPath = componentPath;
        this.width = 0;
        this.height = 0;
        this.rootEle = rootEle;
        this.data = data;
    }
    setTitle(title: string): ModalInstance {
        this._title = title
        return this;
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
        childEle.children[0].setAttribute("ui-id", this.componentPath);
        parserDOM.container.setAttribute("view-path", this.ownerViewPath || this.rootEle.getAttribute("view-path"));
        this._applySize(parserDOM.container)
        const r = document.body.getBoundingClientRect();
        // parserDOM.container.style.left = `${r.width}px`;
        parserDOM.container.style.zIndex = "-1";
        document.body.appendChild(parserDOM.container);
        // this.rootEle.appendChild(parserDOM.container);
        const maxWith = document.body.getBoundingClientRect().width;
        const maxHeight = document.body.getBoundingClientRect().height
        parserDOM.container.style.maxWidth = `${maxWith}px`;
        parserDOM.container.style.maxHeight = `${maxHeight}px`;
        if (!this._title) {
            parserDOM.header.style.display = "none";
        }
        parserDOM.title.innerText = this._title || ' ';
        //document.body.appendChild(parserDOM.container);


        const componentLoader = this.componentModules[`../${this.componentPath}.vue`]
        const Component = (await componentLoader()).default;

        let app = createApp(Component, this.data);
        this.applyVue(app)

        app.mount(parserDOM.body);




        await showElementFromRight(parserDOM.container as HTMLElement, () => {
            parserDOM.container.style.zIndex = "200";
        }, async () => {

        }, 500);

        (new UiWatcher(app, parserDOM.container));
        (new rootWatcher(this.rootEle, parserDOM.container));
        return {
            childEle, app, close() {
                app.unmount();
                parserDOM.container.remove();
            }
        }
    }
    applyVue(app: App<Element>) {
        app.use(PrimeVue, {
            ripple: true,
            locale: {
                startsWith: 'Starts with',
                contains: 'Contains',
                notContains: 'Not contains',
                endsWith: 'Ends with',
                equals: 'Equals',
                notEquals: 'Not equals',
                noFilter: 'No Filter',
                aria: {
                    trueLabel: 'True',
                    falseLabel: 'False'
                },
                firstDayOfWeek: 1,
            },
            theme: {

                options: {
                    darkModeSelector: null,
                    cssLayer: false
                }
            }
        });
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
    load(ownerViewPath: string, rootEle: HTMLElement, componentPath: string, data?: any): ModalInstance {

        let ret = new ModalInstance(ownerViewPath, rootEle, this.htmlModules, this.componentModules, this.htmlLayout, componentPath, data);

        return ret;
    }
}
export default Modal
