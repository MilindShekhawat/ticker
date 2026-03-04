import { ComponentChildren } from "preact";
import { useLocation } from "preact-iso";
import { Button } from "../ui/Button";

type AppFrameProps = {
    section: ComponentChildren;
    children: ComponentChildren;
};

export function AppFrame({ section, children }: AppFrameProps) {
    const { path, route } = useLocation();
    const isSettings = path.startsWith("/dashboard/settings");

    return (
        <div class="min-h-screen bg-zinc-950 text-zinc-100">
            <header class="sticky top-0 z-40 border-b border-zinc-800 bg-zinc-900 backdrop-blur">
                <div class="h-14 w-full grid grid-cols-[112px_1fr_56px_56px] gap-0">
                    <Button
                        class="h-14 w-28 hover:border-zinc-800"
                        onClick={() => route("/dashboard/projects")}
                        title="Projects">
                        <span class="text-lg lowercase font-bold">ticker</span>
                    </Button>

                    <div class="h-14 flex items-center overflow-hidden">{section}</div>

                    <Button
                        variant="primary"
                        size="xs"
                        class="h-14 w-14 border-2"
                        onClick={() => route("/dashboard/settings")}
                        title="Settings">
                        <svg
                            viewBox="0 0 400 400"
                            aria-hidden="true"
                            class="h-8 w-8 fill-current"
                            xmlns="http://www.w3.org/2000/svg">
                            <g>
                                <polygon points="133.3,26.7 160,26.7 160,0 133.3,0 106.7,0 80,0 80,26.7 106.7,26.7" />
                                <polygon points="80,53.3 80,26.7 53.3,26.7 53.3,53.3 26.7,53.3 0,53.3 0,80 26.7,80 53.3,80 53.3,106.7 80,106.7 80,80" />
                                <polygon points="160,53.3 160,80 160,106.7 186.7,106.7 186.7,80 186.7,53.3 186.7,26.7 160,26.7" />
                                <polygon points="373.3,53.3 346.7,53.3 320,53.3 293.3,53.3 266.7,53.3 240,53.3 213.3,53.3 213.3,80 240,80 266.7,80 293.3,80 320,80 346.7,80 373.3,80 400,80 400,53.3" />
                                <polygon points="106.7,106.7 80,106.7 80,133.3 106.7,133.3 133.3,133.3 160,133.3 160,106.7 133.3,106.7" />
                                <polygon points="293.3,160 320,160 320,133.3 293.3,133.3 266.7,133.3 240,133.3 240,160 266.7,160" />
                                <polygon points="133.3,186.7 106.7,186.7 80,186.7 53.3,186.7 26.7,186.7 0,186.7 0,213.3 26.7,213.3 53.3,213.3 80,213.3 106.7,213.3 133.3,213.3 160,213.3 186.7,213.3 186.7,186.7 160,186.7" />
                                <polygon points="240,186.7 240,160 213.3,160 213.3,186.7 213.3,213.3 213.3,240 240,240 240,213.3" />
                                <polygon points="346.7,186.7 346.7,160 320,160 320,186.7 320,213.3 320,240 346.7,240 346.7,213.3 373.3,213.3 400,213.3 400,186.7 373.3,186.7" />
                                <polygon points="266.7,240 240,240 240,266.7 266.7,266.7 293.3,266.7 320,266.7 320,240 293.3,240" />
                                <polygon points="133.3,293.3 160,293.3 160,266.7 133.3,266.7 106.7,266.7 80,266.7 80,293.3 106.7,293.3" />
                                <polygon points="80,320 80,293.3 53.3,293.3 53.3,320 26.7,320 0,320 0,346.7 26.7,346.7 53.3,346.7 53.3,373.3 80,373.3 80,346.7" />
                                <polygon points="160,320 160,346.7 160,373.3 186.7,373.3 186.7,346.7 186.7,320 186.7,293.3 160,293.3" />
                                <polygon points="346.7,320 320,320 293.3,320 266.7,320 240,320 213.3,320 213.3,346.7 240,346.7 266.7,346.7 293.3,346.7 320,346.7 346.7,346.7 373.3,346.7 400,346.7 400,320 373.3,320" />
                                <polygon points="106.7,373.3 80,373.3 80,400 106.7,400 133.3,400 160,400 160,373.3 133.3,373.3" />
                            </g>
                        </svg>
                    </Button>

                    <Button
                        variant="black"
                        size="xs"
                        class="h-14 w-14 hover:border-zinc-950"
                        title="Profile">
                        <img src="/profile.jpg" alt="Profile" class="h-full w-full object-cover" />
                    </Button>
                </div>
            </header>

            <main class="w-full px-4 py-4">{children}</main>
        </div>
    );
}
