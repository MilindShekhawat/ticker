import { useState } from "preact/hooks";
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
import { APIError, api } from "../../lib/api";

type FieldErrors = {
    name?: string;
    email?: string;
    password?: string;
};

function validate(name: string, email: string, password: string): FieldErrors {
    const errors: FieldErrors = {};

    if (!name.trim()) {
        errors.name = "Name is required";
    }

    const normalizedEmail = email.trim().toLowerCase();
    if (!normalizedEmail) {
        errors.email = "Email is required";
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(normalizedEmail)) {
        errors.email = "Enter a valid email";
    }

    if (!password) {
        errors.password = "Password is required";
    } else if (password.length < 10) {
        errors.password = "Password must be at least 10 characters";
    }

    return errors;
}

export function Signup() {
    const { route } = useLocation();

    const [name, setName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    const [formError, setFormError] = useState("");
    const [errors, setErrors] = useState<FieldErrors>({});

    const [isSuccess, setIsSuccess] = useState(false);
    const [isSubmitting, setIsSubmitting] = useState(false);

    const onSubmit = async (event: Event) => {
        event.preventDefault();
        setFormError("");
        setIsSuccess(false);

        const nextErrors = validate(name, email, password);
        setErrors(nextErrors);
        if (Object.keys(nextErrors).length > 0) {
            return;
        }

        setIsSubmitting(true);
        try {
            await api.auth.signup({
                name: name.trim(),
                email: email.trim().toLowerCase(),
                password,
            });

            setIsSuccess(true);
            setPassword("");
            setTimeout(() => route("/"), 900);
        } catch (error) {
            if (error instanceof APIError) {
                setFormError(error.message);
                return;
            }
            setFormError("Could not reach server");
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div class="min-h-screen bg-zinc-950 text-zinc-100 flex items-center justify-center px-4">
            <Card class="w-full max-w-md">
                <CardHeader>
                    <CardTitle class="tracking-tight">Create account</CardTitle>
                    <CardDescription>Team mode signup for Ticker.</CardDescription>
                </CardHeader>
                <CardContent>
                    <form class="flex flex-col gap-4" onSubmit={onSubmit}>
                        <div>
                            <Label htmlFor="name">Name</Label>
                            <Input
                                id="name"
                                type="text"
                                autoComplete="name"
                                value={name}
                                onInput={(e) => setName(e.currentTarget.value)}
                                error={errors.name}
                                disabled={isSubmitting}
                                placeholder="Milind"
                            />
                        </div>
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
                                autoComplete="new-password"
                                value={password}
                                onInput={(e) => setPassword(e.currentTarget.value)}
                                error={errors.password}
                                disabled={isSubmitting}
                                placeholder="At least 10 characters"
                            />
                        </div>

                        {formError && <p class="text-sm text-orange-600">{formError}</p>}
                        {isSuccess && (
                            <p class="text-sm text-lime-400">Signup successful. Redirecting...</p>
                        )}

                        <Button
                            type="submit"
                            variant="primary"
                            class="w-full"
                            disabled={isSubmitting}>
                            {isSubmitting ?
                                <span class="inline-flex items-center gap-2">
                                    creating account
                                    <Loader variant="scanline" />
                                </span>
                            :   "sign up"}
                        </Button>

                        <Button
                            type="button"
                            variant="outline"
                            class="w-full"
                            onClick={() => route("/")}>
                            back to login
                        </Button>
                    </form>
                </CardContent>
            </Card>
        </div>
    );
}
