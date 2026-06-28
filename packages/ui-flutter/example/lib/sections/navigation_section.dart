import 'package:flutter/material.dart';
import 'package:bc_ui_flutter/bc_ui_flutter.dart';
import '../demo_widgets.dart';

/// Navigation widgets: BcBreadcrumb, BcPagination, BcMenubar, BcAffix,
/// BcSearch, BcBackTop.
class NavigationSection extends StatefulWidget {
  const NavigationSection({super.key});

  @override
  State<NavigationSection> createState() => _NavigationSectionState();
}

class _NavigationSectionState extends State<NavigationSection> {
  int _currentPage = 3;
  final ScrollController _scrollController = ScrollController();
  List<BcSearchResult> _searchResults = const [];
  String _searchQuery = '';

  void _onSearch(String query, _) {
    setState(() {
      _searchQuery = query;
      if (query.isEmpty) {
        _searchResults = const [];
        return;
      }
      final all = const [
        BcSearchResult(id: '1', label: 'Dashboard', category: 'Pages'),
        BcSearchResult(id: '2', label: 'Settings', category: 'Pages'),
        BcSearchResult(id: '3', label: 'Users', category: 'Admin'),
        BcSearchResult(id: '4', label: 'Reports', category: 'Analytics'),
        BcSearchResult(id: '5', label: 'Invoices', category: 'Finance'),
        BcSearchResult(id: '6', label: 'Inventory', category: 'Operations'),
      ];
      _searchResults = all
          .where((r) => r.label.toLowerCase().contains(query.toLowerCase()))
          .toList();
    });
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return BcYStack(
      gap: 24,
      children: [
        const SectionHeader(
          title: 'Navigation',
          subtitle: 'Breadcrumbs, pagination, menubar, affix, search, and back-to-top.',
        ),

        // -- BcBreadcrumb --
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

        // -- BcPagination --
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

        // -- BcMenubar --
        DemoCard(
          title: 'BcMenubar',
          description: 'Top menu bar with dropdown menus (File, Edit, View).',
          child: BcMenubar(
            menus: const [
              BcMenubarMenu(label: 'File', items: [
                BcMenubarItem(label: 'New', value: 'new'),
                BcMenubarItem(label: 'Open', value: 'open'),
                BcMenubarItem(label: '', value: '', separator: true),
                BcMenubarItem(label: 'Exit', value: 'exit'),
              ]),
              BcMenubarMenu(label: 'Edit', items: [
                BcMenubarItem(label: 'Undo', value: 'undo'),
                BcMenubarItem(label: 'Redo', value: 'redo'),
                BcMenubarItem(label: '', value: '', separator: true),
                BcMenubarItem(label: 'Cut', value: 'cut'),
                BcMenubarItem(label: 'Copy', value: 'copy'),
                BcMenubarItem(label: 'Paste', value: 'paste'),
              ]),
              BcMenubarMenu(label: 'View', items: [
                BcMenubarItem(label: 'Zoom In', value: 'zoom-in'),
                BcMenubarItem(label: 'Zoom Out', value: 'zoom-out'),
                BcMenubarItem(label: 'Reset', value: 'reset'),
              ]),
            ],
            onSelected: (value, _) {
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text('Menu: $value'),
                  duration: const Duration(milliseconds: 800),
                ),
              );
            },
          ),
        ),

        // -- BcSearch --
        DemoCard(
          title: 'BcSearch',
          description: 'Search input with icon and results dropdown.',
          child: BcYStack(
            gap: 8,
            children: [
              BcSearch(
                value: _searchQuery,
                results: _searchResults,
                placeholder: 'Search pages, settings...',
                emptyText: 'No matches found',
                onSearch: _onSearch,
                onSelect: (result, _) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text('Selected: ${result.label}'),
                      duration: const Duration(milliseconds: 800),
                    ),
                  );
                },
              ),
              if (_searchResults.isNotEmpty)
                Text('${_searchResults.length} result(s)',
                    style: Theme.of(context).textTheme.bodySmall),
            ],
          ),
        ),

        // -- BcAffix + BcBackTop (scrollable demo) --
        DemoCard(
          title: 'BcAffix & BcBackTop',
          description: 'Affix pins a header past a scroll threshold; BackTop shows a floating button after scrolling.',
          child: SizedBox(
            height: 240,
            child: Stack(
              children: [
                Scrollbar(
                  controller: _scrollController,
                  child: ListView.builder(
                    controller: _scrollController,
                    itemCount: 30,
                    itemBuilder: (context, index) {
                      if (index == 0) {
                        return BcAffix(
                          offsetTop: 0,
                          child: Container(
                            color: BcTheme.primary,
                            padding: const EdgeInsets.all(12),
                            child: const Text(
                              'Pinned Header (affix)',
                              style: TextStyle(color: Colors.white, fontWeight: FontWeight.w600),
                            ),
                          ),
                        );
                      }
                      return ListTile(
                        title: Text('Scroll item #$index'),
                      );
                    },
                  ),
                ),
                BcBackTop(
                  visibilityHeight: 100,
                  scrollController: _scrollController,
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }
}