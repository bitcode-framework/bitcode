/// BcApp — root PlatformApp wrapper (ADR-011 Section 3).
///
/// FPW (flutter_platform_widgets) provides [PlatformApp] which auto-selects
/// Material (Android) or Cupertino (iOS/macOS) at runtime based on
/// [PlatformProvider]. BcApp adds BcTheme defaults on top.
///
/// Usage:
/// ```dart
/// void main() => runApp(
///   PlatformProvider(
///     settings: PlatformSettingsData(iosUsesMaterialWidgets: true),
///     child: BcApp(home: MyScreen()),
///   ),
/// );
/// ```

library bc_app;

import 'package:flutter/material.dart';
import 'package:flutter_platform_widgets/flutter_platform_widgets.dart';
import 'theme/bc_theme.dart';

/// Variant override. Defaults to adaptive (platform decides).
enum BcPlatformVariant { adaptive, material, cupertino }

class BcApp extends StatelessWidget {
  final Widget home;
  final BcPlatformVariant variant;
  final String title;
  final bool debugShowCheckedModeBanner;

  const BcApp({
    super.key,
    required this.home,
    this.variant = BcPlatformVariant.adaptive,
    this.title = 'BitCode',
    this.debugShowCheckedModeBanner = false,
  });

  @override
  Widget build(BuildContext context) {
    return PlatformApp(
      title: title,
      debugShowCheckedModeBanner: debugShowCheckedModeBanner,
      material: (context, platform) => MaterialAppData(
        theme: BcTheme.materialLight(),
        darkTheme: BcTheme.materialDark(),
      ),
      cupertino: (context, platform) => CupertinoAppData(
        theme: BcTheme.cupertinoLight(),
      ),
      builder: (context, child) {
        // Force Material variant on non-Apple platforms even when FPW is adaptive.
        if (variant == BcPlatformVariant.material) {
          return Theme(
            data: BcTheme.materialLight(),
            child: Material(child: child),
          );
        }
        return child ?? const SizedBox.shrink();
      },
      home: home,
    );
  }
}
