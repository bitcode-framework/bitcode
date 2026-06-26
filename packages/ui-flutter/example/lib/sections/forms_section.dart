// Forms demo — BcInput, BcSelect, BcCheckbox, BcToggle, BcSegmented, BcSlider.
library forms_section;

import 'package:flutter/material.dart';
import 'package:bc_ui_flutter/bc_ui_flutter.dart';

import '../demo_widgets.dart';

class FormsSection extends StatefulWidget {
  const FormsSection({super.key});

  @override
  State<FormsSection> createState() => _FormsSectionState();
}

class _FormsSectionState extends State<FormsSection> {
  String _inputValue = '';
  String _selectedFruit = 'Apple';
  bool _isChecked = false;
  bool _isToggled = true;
  String _segment = 'Daily';
  double _sliderValue = 50;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SectionHeader(
          title: 'Forms',
          subtitle: 'Adaptive form controls — Material on Android, Cupertino on iOS.',
        ),
        const SizedBox(height: 24),
        DemoCard(
          title: 'BcInput',
          description: 'Adaptive text input (Material TextFormField / CupertinoTextField).',
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              BcInput(
                label: 'Name',
                placeholder: 'Enter your name',
                value: _inputValue,
                onValueChange: (v, _) => setState(() => _inputValue = v),
              ),
              const SizedBox(height: 12),
              BcInput(
                label: 'Email',
                placeholder: 'you@example.com',
                keyboardType: TextInputType.emailAddress,
              ),
              const SizedBox(height: 12),
              BcInput(
                label: 'Disabled',
                placeholder: 'Cannot edit',
                disabled: true,
                defaultValue: 'Read-only value',
              ),
              if (_inputValue.isNotEmpty) ...[
                const SizedBox(height: 8),
                Text('You typed: $_inputValue',
                    style: Theme.of(context).textTheme.bodySmall),
              ],
            ],
          ),
        ),
        const SizedBox(height: 24),
        DemoCard(
          title: 'BcSelect',
          description: 'Adaptive dropdown (Material DropdownButtonFormField / Cupertino modal).',
          child: BcSelect<String>(
            label: 'Favorite Fruit',
            value: _selectedFruit,
            options: const [
              BcSelectOption(value: 'Apple', label: 'Apple'),
              BcSelectOption(value: 'Banana', label: 'Banana'),
              BcSelectOption(value: 'Cherry', label: 'Cherry'),
              BcSelectOption(value: 'Date', label: 'Date'),
            ],
            onValueChange: (v, _) => setState(() => _selectedFruit = v),
          ),
        ),
        const SizedBox(height: 24),
        DemoCard(
          title: 'BcCheckbox & BcToggle',
          description: 'Adaptive checkbox and switch.',
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              BcCheckbox(
                label: 'Accept terms and conditions',
                value: _isChecked,
                onValueChange: (v, _) => setState(() => _isChecked = v),
              ),
              const SizedBox(height: 12),
              BcToggle(
                label: 'Enable notifications',
                value: _isToggled,
                onValueChange: (v, _) => setState(() => _isToggled = v),
              ),
            ],
          ),
        ),
        const SizedBox(height: 24),
        DemoCard(
          title: 'BcSegmented',
          description: 'Adaptive segmented control (Material SegmentedButton / CupertinoSegmentedControl).',
          child: BcSegmented<String>(
            value: _segment,
            options: const [
              BcSegmentedOption(value: 'Daily', label: 'Daily'),
              BcSegmentedOption(value: 'Weekly', label: 'Weekly'),
              BcSegmentedOption(value: 'Monthly', label: 'Monthly'),
            ],
            onValueChange: (v, _) => setState(() => _segment = v),
          ),
        ),
        const SizedBox(height: 24),
        DemoCard(
          title: 'BcSlider',
          description: 'Single-thumb slider with min/max/step clamping.',
          child: BcSlider(
            label: 'Volume',
            value: _sliderValue,
            min: 0,
            max: 100,
            step: 5,
            showValue: true,
            unit: '%',
            onChange: (v, _) => setState(() => _sliderValue = v),
          ),
        ),
      ],
    );
  }
}
