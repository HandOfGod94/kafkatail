require 'avro_turf'

schemas_path = File.join(File.dirname(__FILE__), '..', '..', '..', 'testdata', 'starwars_avro_schemas')
avro = AvroTurf.new(schemas_path: schemas_path)

data = avro.encode({ 'id' => 1,
                     'name' => 'Luke Skywalker',
                     'appearsIn' => %w[EMPIRE JEDI] },
                   schema_name: 'character')

decoded = avro.decode(data, schema_name: 'character')
print(decoded)
