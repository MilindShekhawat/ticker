import { AppFrame } from "../../components/layout/AppFrame";

function Section() {
    return (
        <div class="px-3">
            <p class="text-xs uppercase tracking-wider text-zinc-300">Settings</p>
        </div>
    );
}

export function Settings() {
    return (
        <AppFrame section={<Section />}>
            <div class="border border-zinc-800 bg-zinc-900 p-6 ">
                <h1 class="text-lg font-semibold text-zinc-100 mb-2">Settings</h1>
                <p class="font-mono text-xs text-zinc-400">Settings page coming soon...</p>
            </div>
        </AppFrame>
    );
}
