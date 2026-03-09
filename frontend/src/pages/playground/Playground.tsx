// src/pages/Playground.tsx
import { useState } from "preact/hooks";
import { Button } from "../../components/ui/Button";
import { Input } from "../../components/ui/Input";
import { Textarea } from "../../components/ui/Textarea";
import {
    Card,
    CardHeader,
    CardTitle,
    CardDescription,
    CardContent,
    CardFooter,
} from "../../components/ui/Card";
import { Badge } from "../../components/ui/Badge";
import { Label } from "../../components/ui/Label";
import { Loader } from "../../components/ui/Loader";

const styles = {
    sectionHeading: "text-2xl font-semibold text-zinc-50 mb-4 ",
    cardContentHeading: "text-sm font-medium text-zinc-300 mb-3 ",
    cardContentSection: "flex flex-wrap p-4 gap-3 bg-zinc-900 border border-zinc-700",
};

export function Playground() {
    const [inputValue, setInputValue] = useState("");
    const [textareaValue, setTextareaValue] = useState("");
    const [isLoading, setIsLoading] = useState(false);

    const handleLoadingDemo = () => {
        setIsLoading(true);
        setTimeout(() => setIsLoading(false), 2000);
    };

    return (
        <div class="min-h-screen bg-zinc-950 p-8">
            <div class="max-w-7xl mx-auto space-y-12">
                {/* Header */}
                <div class="text-center">
                    <h1 class="text-4xl font-bold text-zinc-50 mb-2">UI Components Playground</h1>
                    <p class="text-zinc-400">Test all components in dark theme</p>
                </div>

                {/* Cards */}
                <section>
                    <h2 class={`${styles.sectionHeading}`}>CARDS</h2>

                    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {/* Simple Card */}
                        <Card>
                            <CardHeader>
                                <CardTitle>Simple Card</CardTitle>
                                <CardDescription>Basic card with header</CardDescription>
                            </CardHeader>
                            <CardContent>
                                This is a simple card with just header and content.
                            </CardContent>
                        </Card>

                        {/* Card with Footer */}
                        <Card>
                            <CardHeader>
                                <CardTitle>With Footer</CardTitle>
                                <CardDescription>Card with action buttons</CardDescription>
                            </CardHeader>
                            <CardContent>This card has a footer with buttons.</CardContent>
                            <CardFooter class="gap-2">
                                <Button variant="destructive">CANCEL</Button>
                                <Button variant="primary">SAVE</Button>
                            </CardFooter>
                        </Card>

                        {/* Clickable Card */}
                        <Card onClick={() => alert("Card clicked!")}>
                            <CardHeader>
                                <CardTitle>Clickable Card</CardTitle>
                                <CardDescription>Hover to see effect</CardDescription>
                            </CardHeader>
                            <CardContent>This card changes border to lime on hover.</CardContent>
                        </Card>

                        {/* Project Card */}
                        <Card onClick={() => console.log("Project clicked")}>
                            <CardHeader>
                                <CardTitle>Ticker App</CardTitle>
                                <CardDescription>Core product tickets</CardDescription>
                            </CardHeader>
                            <CardContent>
                                <div class="flex items-center gap-4 text-sm text-zinc-400">
                                    <span>12 open</span>
                                    <span>•</span>
                                    <span>45 closed</span>
                                </div>
                            </CardContent>
                        </Card>
                    </div>
                </section>

                {/* Colors */}
                <section>
                    <h2 class={`${styles.sectionHeading}`}>COLORS</h2>
                    <Card>
                        <CardHeader>
                            <CardTitle>COLOR PALETTE</CardTitle>
                            <CardDescription>DESIGN SYSTEM COLORS AND TOKENS</CardDescription>
                        </CardHeader>
                        <CardContent class="space-y-6">
                            {/* Semantic */}
                            <div>
                                <h3 class={`${styles.cardContentHeading}`}>SEMANTIC</h3>
                                <div class="flex flex-wrap gap-3">
                                    {[
                                        { token: "primary", bg: "bg-lime-400" },
                                        { token: "red", bg: "bg-orange-600" },
                                        { token: "green", bg: "bg-green-500" },
                                        { token: "yellow", bg: "bg-yellow-500" },
                                        { token: "blue", bg: "bg-indigo-600" },
                                        { token: "magenta", bg: "bg-pink-500" },
                                        { token: "cyan", bg: "bg-cyan-500" },
                                    ].map(({ token, bg }) => (
                                        <div class="flex flex-col gap-1.5">
                                            <div class={`${bg} w-50 h-40 border border-zinc-700`} />
                                            <span class="text-xs font-bold text-zinc-300 font-mono">
                                                {token}
                                            </span>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </section>

                {/* Buttons */}
                <section>
                    <h2 class={`${styles.sectionHeading}`}>BUTTONS</h2>
                    <Card>
                        <CardHeader>
                            <CardTitle>BUTTON VARIANTS</CardTitle>
                            <CardDescription>ALL BUTTON STYLE AND SIZES</CardDescription>
                        </CardHeader>
                        <CardContent class="space-y-6">
                            {/* Variants */}
                            <div>
                                <h3 class={`${styles.cardContentHeading}`}>VARIANTS</h3>
                                <div class={`${styles.cardContentSection}`}>
                                    <Button variant="primary">Primary</Button>
                                    <Button variant="white">White</Button>
                                    <Button variant="outline">Outline</Button>
                                    <Button variant="black">Black</Button>
                                    <Button variant="invert">Invert</Button>
                                    <Button variant="destructive">Destructive</Button>
                                    <Button disabled>Disabled</Button>
                                </div>
                            </div>

                            {/* Sizes */}
                            <div>
                                <h3 class={`${styles.cardContentHeading}`}>SIZES</h3>
                                <div class={`${styles.cardContentSection}`}>
                                    <Button size="xs" variant="invert">
                                        X-small
                                    </Button>
                                    <Button size="sm" variant="black">
                                        Small
                                    </Button>
                                    <Button size="md" variant="white">
                                        Regular
                                    </Button>
                                    <Button size="lg" variant="outline">
                                        Large
                                    </Button>
                                    <Button size="xl" variant="primary">
                                        X-large
                                    </Button>
                                </div>
                            </div>

                            {/* With Loading */}
                            <div>
                                <h3 class={`${styles.cardContentHeading}`}>LOADING STATE</h3>
                                <div class={`${styles.cardContentSection}`}>
                                    <Button
                                        variant="primary"
                                        onClick={handleLoadingDemo}
                                        disabled={isLoading}>
                                        {isLoading ?
                                            <>
                                                Loading
                                                <Loader variant="scanline" />
                                            </>
                                        :   "Click to load"}
                                    </Button>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </section>

                {/* Inputs */}
                <section>
                    <h2 class={`${styles.sectionHeading}`}>INPUTS</h2>
                    <Card>
                        <CardHeader>
                            <CardTitle>INPUT COMPONENTS</CardTitle>
                            <CardDescription>
                                TEXT COMPONENTS, TEXTAREAS, AND LABELS
                            </CardDescription>
                        </CardHeader>
                        <CardContent class="space-y-6">
                            <div class="space-y-2">
                                <Label htmlFor="email">EMAIL</Label>
                                <Input
                                    id="email"
                                    type="email"
                                    placeholder="you@example.com"
                                    value={inputValue}
                                    onInput={(e) => setInputValue(e.currentTarget.value)}
                                />
                            </div>

                            <div class="space-y-2">
                                <Label htmlFor="password">PASSWORD</Label>
                                <Input id="password" type="password" placeholder="••••••••" />
                            </div>

                            <div class="space-y-2">
                                <Label htmlFor="error-input">INPUT WITH ERROR</Label>
                                <Input
                                    id="error-input"
                                    type="text"
                                    placeholder="This has an error"
                                    error="This field is required"
                                />
                            </div>

                            <div class="space-y-2">
                                <Label htmlFor="disabled">DISABLED INPUT</Label>
                                <Input
                                    id="disabled"
                                    type="text"
                                    placeholder="This is disabled"
                                    disabled
                                />
                            </div>

                            <div class="space-y-2">
                                <Label htmlFor="description">DESCRIPTION</Label>
                                <Textarea
                                    id="description"
                                    placeholder="Enter a description..."
                                    value={textareaValue}
                                    onInput={(e) => setTextareaValue(e.currentTarget.value)}
                                    rows={4}
                                />
                            </div>
                        </CardContent>
                    </Card>
                </section>

                {/* Badges */}
                <section>
                    <h2 class={`${styles.sectionHeading}`}>Badges</h2>
                    <Card>
                        <CardHeader>
                            <CardTitle>BADGE VARIANTS</CardTitle>
                            <CardDescription>STATUS INDICATORS AND TAGS</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <div class="flex flex-wrap gap-2">
                                <Badge variant="primary">Default</Badge>
                                <Badge variant="destructive">Destructive</Badge>
                                <Badge variant="success">Secondary</Badge>
                                <Badge variant="warning">Warning</Badge>
                                <Badge variant="blue">Blue</Badge>
                                <Badge variant="purple">Purple</Badge>
                            </div>

                            <div class="mt-6">
                                <h3 class={`${styles.cardContentHeading}`}>TAG BADGES</h3>
                                <div class="flex flex-wrap gap-2">
                                    <Badge variant="destructive">bug</Badge>
                                    <Badge variant="purple">feature</Badge>
                                    <Badge variant="blue">enhancement</Badge>
                                    <Badge variant="pink">pink</Badge>
                                    <Badge variant="warning">urgent</Badge>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </section>

                {/* Loading States */}
                <section>
                    <h2 class={`${styles.sectionHeading}`}>Loading States</h2>
                    <Card>
                        <CardHeader>
                            <CardTitle>LOADERS</CardTitle>
                            <CardDescription>Loading indicators</CardDescription>
                        </CardHeader>
                        <CardContent class="space-y-6">
                            <div>
                                <h3 class={`${styles.cardContentHeading}`}>NOISE</h3>
                                <div class={`${styles.cardContentSection}`}>
                                    <Loader />
                                </div>
                            </div>
                            <div>
                                <h3 class={`${styles.cardContentHeading}`}>PULSE</h3>
                                <div class={`${styles.cardContentSection}`}>
                                    <Loader variant="pulse" />
                                </div>
                            </div>
                            <div>
                                <h3 class={`${styles.cardContentHeading}`}>BUTTON</h3>
                                <div class={`${styles.cardContentSection}`}>
                                    <Button>
                                        Loading
                                        <Loader variant="scanline" />
                                    </Button>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </section>
            </div>
        </div>
    );
}
