import 'package:flutter_test/flutter_test.dart';
import 'package:bc_ui_flutter/src/core/types.dart';

void main() {
  group('BcChangeMeta', () {
    test('default source is programmatic', () {
      const meta = BcChangeMeta(reason: BcChangeReason.input);
      expect(meta.source, BcEventSource.programmatic);
    });

    test('equality based on reason and source', () {
      const a = BcChangeMeta(reason: BcChangeReason.select, source: BcEventSource.pointer);
      const b = BcChangeMeta(reason: BcChangeReason.select, source: BcEventSource.pointer);
      const c = BcChangeMeta(reason: BcChangeReason.toggle, source: BcEventSource.pointer);
      expect(a == b, isTrue);
      expect(a == c, isFalse);
      expect(a.hashCode, b.hashCode);
    });

    test('nativeEvent does not affect equality', () {
      const a = BcChangeMeta(reason: BcChangeReason.input, nativeEvent: 'event1');
      const b = BcChangeMeta(reason: BcChangeReason.input, nativeEvent: 'event2');
      expect(a == b, isTrue);
    });
  });

  group('BcEventSource', () {
    test('has all expected values', () {
      expect(BcEventSource.values.length, 4);
      expect(BcEventSource.values, contains(BcEventSource.pointer));
      expect(BcEventSource.values, contains(BcEventSource.keyboard));
      expect(BcEventSource.values, contains(BcEventSource.programmatic));
      expect(BcEventSource.values, contains(BcEventSource.external));
    });
  });

  group('BcChangeReason', () {
    test('has all expected values', () {
      expect(BcChangeReason.values.length, 6);
      expect(BcChangeReason.values, contains(BcChangeReason.input));
      expect(BcChangeReason.values, contains(BcChangeReason.reaction));
      expect(BcChangeReason.values, contains(BcChangeReason.programmatic));
      expect(BcChangeReason.values, contains(BcChangeReason.select));
      expect(BcChangeReason.values, contains(BcChangeReason.toggle));
      expect(BcChangeReason.values, contains(BcChangeReason.submit));
    });
  });

  group('BcComponentSize', () {
    test('has sm, md, lg', () {
      expect(BcComponentSize.values.length, 3);
    });
  });

  group('BcComponentVariant', () {
    test('has all expected variants', () {
      expect(BcComponentVariant.values.length, 7);
    });
  });
}
