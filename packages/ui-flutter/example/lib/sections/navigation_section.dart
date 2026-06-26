import 'package:flutter/material.dart';
import 'package:bc_ui_flutter/bc_ui_flutter.dart';
import '../demo_widgets.dart';

/// Navigation widgets: BcBreadcrumb, BcPagination.
class NavigationSection extends StatefulWidget {
  const NavigationSection({super.key});

  @override
  State<NavigationSection> createState() => _NavigationSectionState();
}

class _NavigationSectionState extends State<NavigationSection> {
  int _currentPage = 3;

  @override
  Widget build(BuildContext context) {
    return BcYStack(
      gap: 24,
      children: [
        const SectionHeader(
          title: 'Navigation',
          subtitle: 'Breadcrumbs and pagination controls.',
        ),

        // ── BcBreadcrumb ──
        DemoCard(
          title: 'BcBreadcrumb',
          description: 'Adaptive breadcrumb navigation trail with separators.',
          child: BcBreadcrumb(
            items: const [
              BcBreadcrumbItem(label: 'Home', path: '/'),
              BcBreadcrumbItem(label: 'Products', path: '/products'),
              BcBreadcrumbItem(
                  label: 'Electronics', path: '/products/electronics'),
              BcBreadcrumbItem(label: 'Laptops', active: true),
            ],
            onSelected: (path, _) {
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text('Navigate to: $path'),
                  duration: const Duration(milliseconds: 800),
                ),
              );
            },
          ),
        ),

        // ── BcPagination ──
        DemoCard(
          title: 'BcPagination',
          description: 'Page navigation with sibling range and ellipsis.',
          child: BcYStack(
            gap: 12,
            children: [
              BcPagination(
                currentPage: _currentPage,
                totalPages: 10,
                showFirstLast: true,
                showPrevNext: true,
                onPageChange: (page, _) => setState(() => _currentPage = page),
              ),
              Text('Current page: $_currentPage',
                  style: Theme.of(context).textTheme.bodySmall),
            ],
          ),
        ),
      ],
    );
  }
}
