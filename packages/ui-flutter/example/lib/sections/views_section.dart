// Views demo - BcViewMap, BcViewReport, BcViewEditor, BcViewActivity,
// BcViewGantt, BcKanbanBoard, BcVirtualScroller, BcTieredMenu, BcOrgChart,
// BcCarousel.
library views_section;

import 'package:flutter/material.dart';
import 'package:bc_ui_flutter/bc_ui_flutter.dart';

import '../demo_widgets.dart';

class ViewsSection extends StatelessWidget {
  const ViewsSection({super.key});

  @override
  Widget build(BuildContext context) {
    return BcYStack(
      gap: 24,
      children: [
        const SectionHeader(
          title: 'Views',
          subtitle: 'Enterprise data views: Kanban, Gantt, maps, reports, activity, org charts, and more.',
        ),

        // -- BcKanbanBoard --
        DemoCard(
          title: 'BcKanbanBoard',
          description: 'Kanban board with cross-column drag-drop.',
          child: SizedBox(
            height: 320,
            child: SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              child: SizedBox(
                width: 720,
                child: BcKanbanBoard(
                  boardTitle: 'Sprint Board',
                  columns: const [
                    BcKanbanColumnConfig(id: 'todo', name: 'To Do'),
                    BcKanbanColumnConfig(id: 'doing', name: 'In Progress'),
                    BcKanbanColumnConfig(id: 'done', name: 'Done'),
                  ],
                  cards: const [
                    BcKanbanCard(id: '1', title: 'Setup CI', column: 'todo'),
                    BcKanbanCard(id: '2', title: 'Fix login bug', column: 'doing'),
                    BcKanbanCard(id: '3', title: 'Write docs', column: 'doing'),
                    BcKanbanCard(id: '4', title: 'Design homepage', column: 'done'),
                  ],
                ),
              ),
            ),
          ),
        ),

        // -- BcViewGantt --
        DemoCard(
          title: 'BcViewGantt',
          description: 'Gantt chart with task bars and timeline.',
          child: SizedBox(
            height: 260,
            child: BcViewGantt(
              tasks: const [
                BcGanttTask(id: 't1', label: 'Planning', start: 0, duration: 5, progress: 100),
                BcGanttTask(id: 't2', label: 'Design', start: 3, duration: 7, progress: 60),
                BcGanttTask(id: 't3', label: 'Development', start: 8, duration: 12, progress: 20),
                BcGanttTask(id: 't4', label: 'Testing', start: 18, duration: 5, progress: 0),
              ],
              timelineStart: '2026-06-01',
              timelineEnd: '2026-07-15',
              zoomLevel: BcGanttZoomLevel.week,
              showProgress: true,
              showTodayLine: true,
            ),
          ),
        ),

        // -- BcViewMap --
        DemoCard(
          title: 'BcViewMap',
          description: 'Map view with geo markers and tap selection.',
          child: BcViewMap(
            height: 240,
            center: const BcGeoCoord(lat: -6.2088, lng: 106.8456),
            zoom: 10,
            markers: const [
              BcMapMarker(id: 'hq', lat: -6.2088, lng: 106.8456, label: 'HQ'),
              BcMapMarker(id: 'branch', lat: -6.2200, lng: 106.8600, label: 'Branch Office'),
              BcMapMarker(id: 'warehouse', lat: -6.1900, lng: 106.8300, label: 'Warehouse'),
            ],
          ),
        ),

        // -- BcViewReport --
        DemoCard(
          title: 'BcViewReport',
          description: 'Structured report with titled sections.',
          child: BcViewReport(
            title: 'Q2 2026 Summary',
            subtitle: 'Engineering Department',
            sections: const [
              BcReportSection(
                id: 's1',
                title: 'Highlights',
                content: 'Shipped 14 demo sections, reached full component parity.',
              ),
              BcReportSection(
                id: 's2',
                title: 'Metrics',
                content: '11 widget tests pass, 0 analyze issues, 95+ components demoed.',
              ),
              BcReportSection(
                id: 's3',
                title: 'Next Quarter',
                content: 'Focus on BcRenderer, reactive engine wiring, and offline sync.',
              ),
            ],
          ),
        ),

        // -- BcViewActivity --
        DemoCard(
          title: 'BcViewActivity',
          description: 'Activity feed timeline sorted newest-first.',
          child: SizedBox(
            height: 280,
            child: BcViewActivity(
              entries: [
                BcViewActivityEntry(
                  id: '1',
                  timestamp: DateTime(2026, 6, 28, 14, 30),
                  actor: 'Alice',
                  action: 'merged PR #42',
                  detail: 'Flutter demo parity complete',
                ),
                BcViewActivityEntry(
                  id: '2',
                  timestamp: DateTime(2026, 6, 28, 10, 15),
                  actor: 'Bob',
                  action: 'opened issue #99',
                  detail: 'BcChart crash in headless tests',
                ),
                BcViewActivityEntry(
                  id: '3',
                  timestamp: DateTime(2026, 6, 27, 16, 0),
                  actor: 'Charlie',
                  action: 'pushed 3 commits',
                ),
              ],
            ),
          ),
        ),

        // -- BcViewEditor --
        DemoCard(
          title: 'BcViewEditor',
          description: 'Visual layout editor with drag-to-reorder components (disabled here to show static layout).',
          child: BcViewEditor(
            disabled: true,
            defaultComponents: const [
              BcViewEditorComponent(id: 'c1', type: 'list', label: 'Task List'),
              BcViewEditorComponent(id: 'c2', type: 'chart', label: 'Burndown Chart'),
              BcViewEditorComponent(id: 'c3', type: 'form', label: 'New Task Form'),
            ],
          ),
        ),

        // -- BcOrgChart --
        DemoCard(
          title: 'BcOrgChart',
          description: 'Organizational hierarchy with expandable nodes.',
          child: BcOrgChart(
            defaultExpanded: const ['ceo', 'vp1'],
            root: const BcOrgNode(
              id: 'ceo',
              label: 'CEO',
              title: 'Chief Executive',
              children: [
                BcOrgNode(
                  id: 'vp1',
                  label: 'VP Engineering',
                  title: 'Engineering',
                  children: [
                    BcOrgNode(id: 'lead1', label: 'Tech Lead', title: 'Frontend'),
                    BcOrgNode(id: 'lead2', label: 'Tech Lead', title: 'Backend'),
                  ],
                ),
                BcOrgNode(
                  id: 'vp2',
                  label: 'VP Sales',
                  title: 'Sales',
                ),
              ],
            ),
          ),
        ),

        // -- BcTieredMenu --
        DemoCard(
          title: 'BcTieredMenu',
          description: 'Cascading tiered menu with submenus.',
          child: BcTieredMenu(
            items: const [
              BcTieredMenuItem(
                key: 'new',
                label: 'New',
                icon: Icons.add,
                children: [
                  BcTieredMenuItem(key: 'project', label: 'Project'),
                  BcTieredMenuItem(key: 'file', label: 'File'),
                  BcTieredMenuItem(key: 'folder', label: 'Folder'),
                ],
              ),
              BcTieredMenuItem(
                key: 'export',
                label: 'Export As',
                icon: Icons.download,
                children: [
                  BcTieredMenuItem(key: 'csv', label: 'CSV'),
                  BcTieredMenuItem(key: 'pdf', label: 'PDF'),
                ],
              ),
              BcTieredMenuItem(key: 'exit', label: 'Exit', icon: Icons.exit_to_app),
            ],
          ),
        ),

        // -- BcVirtualScroller --
        DemoCard(
          title: 'BcVirtualScroller',
          description: 'Windowed virtual list for large datasets.',
          child: SizedBox(
            height: 200,
            child: BcVirtualScroller(
              items: List.generate(
                500,
                (i) => BcVirtualScrollerItem(
                  id: '$i',
                  label: 'Item ${i + 1}',
                  description: 'Description for item ${i + 1}',
                ),
              ),
            ),
          ),
        ),

        // -- BcCarousel --
        DemoCard(
          title: 'BcCarousel',
          description: 'Image carousel with autoplay, dots, and arrows.',
          child: BcCarousel(
            aspectRatio: 16 / 9,
            autoplay: false,
            slides: const [
              BcCarouselSlide(caption: 'Slide One'),
              BcCarouselSlide(caption: 'Slide Two'),
              BcCarouselSlide(caption: 'Slide Three'),
            ],
          ),
        ),
      ],
    );
  }
}
