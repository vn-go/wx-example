import libs from '@core/lib';

class Views extends libs.BaseUI {
    onInit() {
        this.loadAllViews();
    }
    async loadAllViews() {

    }
    showViewPath() {
        alert(this.getViewPath());
    }
}

export default Views;