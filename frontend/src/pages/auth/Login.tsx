import { useEffect, useRef, useState } from "preact/hooks";
import { useLocation } from "preact-iso";
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "../../components/ui/Card";
import { Button } from "../../components/ui/Button";
import { Input } from "../../components/ui/Input";
import { Label } from "../../components/ui/Label";
import { Loader } from "../../components/ui/Loader";
import { APIError, api, isAbortError } from "../../lib/api";

type FieldErrors = {
    email?: string;
    password?: string;
};

function validate(email: string, password: string): FieldErrors {
    const errors: FieldErrors = {};
    const normalizedEmail = email.trim().toLowerCase();

    if (!normalizedEmail) {
        errors.email = "Email is required";
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(normalizedEmail)) {
        errors.email = "Enter a valid email";
    }

    if (!password) {
        errors.password = "Password is required";
    }

    return errors;
}

export function Login() {
    const { route } = useLocation();
    const requestControllerRef = useRef<AbortController | null>(null);

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    const [formError, setFormError] = useState("");
    const [errors, setErrors] = useState<FieldErrors>({});

    const [isSubmitting, setIsSubmitting] = useState(false);

    const onSubmit = async (event: Event) => {
        event.preventDefault();
        setFormError("");

        const normalizedEmail = email.trim().toLowerCase();

        const nextErrors = validate(normalizedEmail, password);
        setErrors(nextErrors);
        if (Object.keys(nextErrors).length > 0) return;

        requestControllerRef.current?.abort();
        const controller = new AbortController();
        requestControllerRef.current = controller;

        setIsSubmitting(true);
        try {
            await api.auth.login(
                { email: normalizedEmail, password },
                { signal: controller.signal },
            );

            route("/dashboard/projects");
        } catch (error) {
            if (isAbortError(error)) return;

            if (error instanceof APIError) {
                setFormError(error.message);
                return;
            }

            setFormError("Could not reach server");
        } finally {
            setIsSubmitting(false);
        }
    };

    useEffect(() => {
        return () => requestControllerRef.current?.abort();
    }, []);

    return (
        <div class="min-h-screen bg-zinc-950 text-zinc-100 flex items-center justify-center px-4">
            <Card class="w-full max-w-md">
                <CardHeader>
                    <CardTitle class="tracking-tight">Login to Ticker</CardTitle>
                    <CardDescription>Team mode login for Ticker.</CardDescription>
                </CardHeader>
                <CardContent>
                    <form class="flex flex-col gap-4" onSubmit={onSubmit}>
                        <div>
                            <Label htmlFor="email">Email</Label>
                            <Input
                                id="email"
                                type="email"
                                autoComplete="email"
                                value={email}
                                onInput={(e) => setEmail(e.currentTarget.value)}
                                error={errors.email}
                                disabled={isSubmitting}
                                placeholder="you@company.com"
                            />
                        </div>
                        <div>
                            <Label htmlFor="password">Password</Label>
                            <Input
                                id="password"
                                type="password"
                                autoComplete="current-password"
                                value={password}
                                onInput={(e) => setPassword(e.currentTarget.value)}
                                error={errors.password}
                                disabled={isSubmitting}
                                placeholder="Your password"
                            />
                        </div>

                        {formError && <p class="text-sm text-orange-600">{formError}</p>}

                        <Button
                            type="submit"
                            variant="primary"
                            class="w-full"
                            disabled={isSubmitting}>
                            {isSubmitting ?
                                <>
                                    logging in
                                    <Loader variant="scanline" />
                                </>
                            :   "login"}
                        </Button>

                        <Button
                            type="button"
                            variant="outline"
                            class="w-full"
                            onClick={() => route("/signup")}>
                            create account
                        </Button>
                    </form>
                </CardContent>
            </Card>
        </div>
    );
}
