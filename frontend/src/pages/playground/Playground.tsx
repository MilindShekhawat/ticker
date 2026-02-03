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
                                    <Button variant="primary">PRIMARY</Button>
                                    <Button variant="white">WHITE</Button>
                                    <Button variant="outline">OUTLINE</Button>
                                    <Button variant="black">BLACK</Button>
                                    <Button variant="invert">INVERT</Button>
                                    <Button variant="destructive">DESTRUCTIVE</Button>
                                    <Button disabled>DISABLED</Button>
                                </div>
                            </div>

                            {/* Sizes */}
                            <div>
                                <h3 class={`${styles.cardContentHeading}`}>SIZES</h3>
                                <div class={`${styles.cardContentSection}`}>
                                    <Button size="xs" variant="invert">
                                        X-SMALL
                                    </Button>
                                    <Button size="sm" variant="black">
                                        SMALL
                                    </Button>
                                    <Button size="md" variant="white">
                                        REGULAR
                                    </Button>
                                    <Button size="lg" variant="outline">
                                        LARGE
                                    </Button>
                                    <Button size="xl" variant="primary">
                                        X-LARGE
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
                                            <span class="flex items-center gap-2">
                                                <Loader />
                                                LOADING...
                                            </span>
                                        :   "CLICK TO LOAD"}
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
                                <Badge variant="success">Secondary</Badge>
                                <Badge variant="destructive">Destructive</Badge>
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
                                    <Loader variant="dense"/>
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
