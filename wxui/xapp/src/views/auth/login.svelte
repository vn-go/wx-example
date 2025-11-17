<script lang="ts">
    import { BaseUi } from "@core/core";
    import LoginLayout from "@template-layout/LoginLayout.svelte";
    const test = () => {
        alert("OK");
    };
    class Login extends BaseUi {
        data = $state({
            username: "",
            password: "",
            error: "",
        });
        _afterLogin: (() => void) | undefined;
        constructor() {
            super();
        }
        async doLogin() {
            try {
                let ret = await this.app.loginAsync(
                    this.data.username,
                    this.data.password,
                );

                if (ret.error) {
                    this.data.error = "Login fail";
                } else {
                    this.app.setSessionValue("tk", ret.access_token);
                    if (this._afterLogin) {
                        this._afterLogin();
                    }
                }
            } catch (error) {
                this.data.error = "Login fail";
            }
        }
        afterLogin(fn: () => void) {
            this._afterLogin = fn;
        }
    }
    const login = new Login(); // lam sao tra ve cai nay sau khi goi ham renderDynamicComponent
</script>

<div bind:this={login.Element} class="flex-1 w-full h-full overflow-hidden">
    <LoginLayout>
        <input
            slot="username-input"
            bind:value={login.data.username}
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500"
        />
        <input
            type="password"
            slot="password-input"
            bind:value={login.data.password}
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500"
        />
        <button
            slot="submit-button"
            on:click={() => login.doLogin()}
            class="w-full py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-xl font-semibold hover:shadow-lg transform hover:scale-105 transition-all duration-200 flex items-center justify-center gap-2"
        >
            Sign In
        </button>
        <div
            slot="error-message"
            class="text-red-600 text-sm font-medium animate-fade-in"
        >
            {login.data.error}
        </div>
    </LoginLayout>
</div>
