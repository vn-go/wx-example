<script lang="ts">
  import SideBar from "@components/sidbar.svelte";
  import { application, BaseUi } from "@core/core";
  import Layout from "@template-layout/app.svelte";
  application.onUrlChange((data) => {
    console.log(data.PagePath);
    console.log(data.Query);
  });

  class AppUI extends BaseUi {
    bodyEl: HTMLElement | undefined;
    checkAuth() {
      let token = this.app.getSessionValue<string>("token");
      if (!token) {
      }
    }
    onInit(): void {
      this.app.onUrlChange(async (data) => {
        let token = this.app.getSessionValue<string>("tk");
        if (!token) {
          let ret = await this.loadLogin(data.PagePath);
        } else {
          if (data.PagePath) {
            let ret = await this.app.renderAsync(data.PagePath, this.bodyEl);
            if (!ret.success) {
              let retErrPage = await this.app.loadErrorPage(this.bodyEl);
              await this.app.slideReplaceRight(
                (retErrPage as any).ele,
                this.bodyEl,
              );
            } else {
              await this.app.slideReplaceRight((ret as any).ele, this.bodyEl);
            }
          }
        }
      });
    }
    async loadLogin(retPage: string | undefined) {
      let ret = await this.app.renderAsync("auth/login", this.bodyEl);
      if (!ret.success) {
        let ret = await this.app.loadErrorPage(this.bodyEl);
        await this.app.slideReplaceRight((ret as any).ele, this.bodyEl);
      } else {
        (ret as any).ins.afterLogin(async () => {
          if (retPage) {
            try {
              let retCmp = await this.app.renderAsync(retPage, this.bodyEl);
              if (!retCmp.success) {
                let retErrPage = await this.app.loadErrorPage(this.bodyEl);
                this.app.slideReplaceRight(
                  (retErrPage as any).ele,
                  this.bodyEl,
                );
              } else {
                this.app.slideReplaceRight((retCmp as any).ele, this.bodyEl);
              }
            } catch (error) {
              let retErrPage = await this.app.loadErrorPage(this.bodyEl);
              await this.app.slideReplaceRight(
                (retErrPage as any).ele,
                this.bodyEl,
              );
            }
          }
        });
      }
    }
  }
  let container: HTMLDivElement;
  const appUi = new AppUI();
</script>

<main>
  <Layout>
    <div slot="sidebar">
      <SideBar />
    </div>

    <div slot="body" bind:this={appUi.bodyEl}></div>
  </Layout>
</main>
