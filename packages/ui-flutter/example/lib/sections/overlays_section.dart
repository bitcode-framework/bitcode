// Overlays demo — BcDialog, BcAlertDialog, BcTooltip, BcAccordion,
// BcDisclosure, BcTabs.
library overlays_section;

import 'package:flutter/material.dart';
import 'package:flutter_platform_widgets/flutter_platform_widgets.dart';
import 'package:bc_ui_flutter/bc_ui_flutter.dart';

import '../demo_widgets.dart';

class OverlaysSection extends StatelessWidget {
  const OverlaysSection({super.key});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SectionHeader(
          title: 'Overlays',
          subtitle: 'Dialogs, tooltips, accordion, disclosure, and tabs.',
        ),
        const SizedBox(height: 24),
        DemoCard(
          title: 'BcDialog',
          description: 'Adaptive confirm/cancel dialog (Material AlertDialog / CupertinoAlertDialog).',
          child: PlatformElevatedButton(
            onPressed: () async {
              final result = await BcDialog.show(
                context: context,
                title: 'Delete Item',
                content: const Text(
                    'Are you sure you want to delete this item? This action cannot be undone.'),
                confirmText: 'Delete',
                cancelText: 'Cancel',
              );
              if (context.mounted) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(result == true ? 'Confirmed' : 'Cancelled'),
                    duration: const Duration(seconds: 1),
                  ),
                );
              }
            },
            child: const Text('Open Dialog'),
          ),
        ),
        const SizedBox(height: 24),
        DemoCard(
          title: 'BcAlertDialog',
          description: 'Simple alert with a single dismiss button.',
          child: PlatformElevatedButton(
            onPressed: () {
              BcAlertDialog.show(
                context: context,
                title: 'Information',
                content: const Text(
                    'This is an informational alert dialog. Tap OK to dismiss.'),
                dismissText: 'OK',
              );
            },
            child: const Text('Open Alert'),
          ),
        ),
        const SizedBox(height: 24),
        DemoCard(
          title: 'BcTooltip',
          description: 'Adaptive tooltip on long-press/hover.',
          child: BcTooltip(
            message: 'This is a helpful tooltip!',
            child: Icon(
              PlatformIcons(context).info,
              color: BcTheme.primary,
              size: 32,
            ),
          ),
        ),
        const SizedBox(height: 24),
        DemoCard(
          title: 'BcAccordion',
          description: 'Collapsible accordion sections (Material ExpansionTile).',
          child: BcAccordion(
            initialExpandedIndex: 0,
            items: const [
              BcAccordionItem(
                title: 'What is BitCode UI?',
                child: Text(
                    'BitCode UI is a cross-framework component library that generates React, Vue, Stencil, and Flutter widgets from a single Mitosis source.'),
              ),
              BcAccordionItem(
                title: 'How does it work?',
                child: Text(
                    'Mitosis compiles .lite.tsx files into target frameworks. The Flutter adapter ports the contracts to Dart widgets.'),
              ),
              BcAccordionItem(
                title: 'Is it production ready?',
                child: Text(
                    'The Flutter adapter has 1500+ tests passing and clean static analysis. It is ready for integration.'),
              ),
            ],
          ),
        ),
        const SizedBox(height: 24),
        DemoCard(
          title: 'BcDisclosure',
          description: 'Collapsible panel with animated expand/collapse.',
          child: BcDisclosure(
            title: 'Advanced Settings',
            defaultOpen: false,
            child: const Padding(
              padding: EdgeInsets.only(top: 8),
              child: Text(
                'These are advanced settings that were hidden behind a '
                'disclosure panel. Tap the header to toggle.',
                style: TextStyle(fontSize: 14),
              ),
            ),
          ),
        ),
        const SizedBox(height: 24),
        DemoCard(
          title: 'BcTabs',
          description: 'Material 3 tab bar with swipeable content views.',
          child: SizedBox(
            height: 200,
            child: BcTabs(
              tabs: const [
                BcTabItem(
                  label: 'Overview',
                  child: Center(child: Text('Overview tab content')),
                ),
                BcTabItem(
                  label: 'Details',
                  child: Center(child: Text('Details tab content')),
                ),
                BcTabItem(
                  label: 'Reviews',
                  child: Center(child: Text('Reviews tab content')),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }
}
