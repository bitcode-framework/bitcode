// Overlays demo - BcDialog, BcAlertDialog, BcTooltip, BcAccordion,
// BcDisclosure, BcTabs, BcSplitter, BcDropdown, BcPopover, BcCommand,
// BcHoverCard, BcActionsheet.
library overlays_section;

import 'package:flutter/material.dart';
import 'package:flutter_platform_widgets/flutter_platform_widgets.dart';
import 'package:bc_ui_flutter/bc_ui_flutter.dart';

import '../demo_widgets.dart';

class OverlaysSection extends StatefulWidget {
  const OverlaysSection({super.key});

  @override
  State<OverlaysSection> createState() => _OverlaysSectionState();
}

class _OverlaysSectionState extends State<OverlaysSection> {
  String _dropdownValue = 'Option A';
  bool _actionsheetOpen = false;
  String _lastAction = 'none';

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SectionHeader(
          title: 'Overlays',
          subtitle: 'Dialogs, tooltips, dropdowns, splitters, command palette, hover cards, and action sheets.',
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
        const SizedBox(height: 24),
        DemoCard(
          title: 'BcSplitter',
          description: 'Draggable two-pane splitter (horizontal or vertical).',
          child: SizedBox(
            height: 160,
            child: BcSplitter(
              axis: Axis.horizontal,
              initialFraction: 0.5,
              first: Container(
                color: BcTheme.primary.withValues(alpha: 0.15),
                alignment: Alignment.center,
                child: const Text('Pane 1'),
              ),
              second: Container(
                color: BcTheme.secondary.withValues(alpha: 0.15),
                alignment: Alignment.center,
                child: const Text('Pane 2'),
              ),
            ),
          ),
        ),
        const SizedBox(height: 24),
        DemoCard(
          title: 'BcDropdown',
          description: 'Adaptive dropdown menu with a trigger button.',
          child: BcDropdown<String>(
            value: _dropdownValue,
            options: const [
              BcSelectOption(value: 'Option A', label: 'Option A'),
              BcSelectOption(value: 'Option B', label: 'Option B'),
              BcSelectOption(value: 'Option C', label: 'Option C'),
            ],
            placeholder: 'Choose...',
            onValueChange: (v, _) => setState(() => _dropdownValue = v),
          ),
        ),
        const SizedBox(height: 24),
        DemoCard(
          title: 'BcPopover',
          description: 'Overlay popover anchored to a trigger widget (Material-only).',
          child: BcPopover(
            title: 'Popover Title',
            width: 220,
            trigger: ElevatedButton.icon(
              onPressed: null,
              icon: const Icon(Icons.touch_app_outlined),
              label: const Text('Tap to open'),
            ),
            child: const Padding(
              padding: EdgeInsets.all(8),
              child: Text(
                'This popover content appears in an overlay near the trigger. '
                'Tap outside to dismiss.',
                style: TextStyle(fontSize: 13),
              ),
            ),
          ),
        ),
        const SizedBox(height: 24),
        DemoCard(
          title: 'BcCommand',
          description: 'Command palette triggered by Ctrl+K / Cmd+K (Material-only).',
          child: BcCommand(
            placeholder: 'Type a command...',
            commands: const [
              BcCommandItem(key: 'new-file', label: 'New File'),
              BcCommandItem(key: 'open-file', label: 'Open File'),
              BcCommandItem(key: 'save', label: 'Save'),
              BcCommandItem(key: 'print', label: 'Print'),
              BcCommandItem(key: 'exit', label: 'Exit'),
            ],
            child: const Text(
              'Focus this area and press Ctrl+K (or Cmd+K) to open the command palette.',
              style: TextStyle(fontSize: 13),
            ),
          ),
        ),
        const SizedBox(height: 24),
        DemoCard(
          title: 'BcHoverCard',
          description: 'Card that reveals on hover with placement control.',
          child: BcHoverCard(
            placement: BcHoverCardPlacement.bottom,
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
              decoration: BoxDecoration(
                color: BcTheme.primary,
                borderRadius: BorderRadius.circular(6),
              ),
              child: const Text(
                'Hover me',
                style: TextStyle(color: Colors.white),
              ),
            ),
            content: Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: BcTheme.card,
                border: Border.all(color: BcTheme.border),
              ),
              child: const Text(
                'This card appears on hover and stays while the pointer '
                'remains over the trigger or the card.',
                style: TextStyle(fontSize: 13),
              ),
            ),
          ),
        ),
        const SizedBox(height: 24),
        DemoCard(
          title: 'BcActionsheet',
          description: 'Bottom sheet action menu that slides up from the bottom.',
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              ElevatedButton(
                onPressed: () => setState(() => _actionsheetOpen = true),
                child: const Text('Open Action Sheet'),
              ),
              const SizedBox(height: 8),
              Text('Last action: $_lastAction',
                  style: Theme.of(context).textTheme.bodySmall),
              BcActionsheet(
                open: _actionsheetOpen,
                title: 'Share',
                description: 'Choose how to share this item.',
                items: const [
                  BcActionsheetItem(
                      label: 'Copy Link', value: 'copy', icon: Icons.link),
                  BcActionsheetItem(
                      label: 'Email', value: 'email', icon: Icons.email_outlined),
                  BcActionsheetItem(
                      label: 'Delete',
                      value: 'delete',
                      icon: Icons.delete_outline,
                      destructive: true),
                ],
                onSelect: (item, _) =>
                    setState(() {
                  _actionsheetOpen = false;
                  _lastAction = item.label;
                }),
                onCancel: (_, _) => setState(() => _actionsheetOpen = false),
              ),
            ],
          ),
        ),
      ],
    );
  }
}
