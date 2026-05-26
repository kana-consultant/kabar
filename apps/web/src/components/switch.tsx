import { Switch } from "@kana-consultant/ui-kit"
import { setThemeMode, useResolvedTheme } from "@kana-consultant/ui-kit"

export function ThemeSwitch() {
  const resolvedTheme = useResolvedTheme()

  return (
    <Switch
      checked={resolvedTheme === "dark"}
      onCheckedChange={(checked : any) => {
        if (checked) {
          setThemeMode("dark")
        } else {
          setThemeMode("light")
        }
      }}
    />
  )
}