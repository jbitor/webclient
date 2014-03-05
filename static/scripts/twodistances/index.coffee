# All distances must be between 0.0 and 1.0.

exports.Coordinates =
class Coordinates
  constructor: (@x, @y) ->

  toString: ->
    "Coordinates(#{@x}, #{@y})"

exports.Node =
class Node
  constructor: (@distanceA, @distanceB, @source, @color) ->
    if not (@distanceA >=0 && @distanceA <= 1)
      throw new Error("distanceA=#{@distanceA}", @distanceA)

    if not (@distanceB >=0 && @distanceB <= 1)
      throw new Error("distanceB=#{@distanceB}", @distanceB)

    if not @source?
      @source = null
    else if @source not instanceof Node
      throw new Error("source=#{@source}", @source)

    if not @color?
      @color = null
    else if typeof @color != 'string'
      throw new Error("color=#{@color}")

    @twoDistances = null

  setTwoDistances: (@twoDistances) ->
    # it is an error if the distances sum to less than the direct origin distances
    if (@distanceA + @distanceB < @twoDistances.originDistances)
      console.error "Distances are impossibly small.", @distanceA, @distanceB
      throw new Error("Distances are impossibly small.", this)


    # or if difference between the distances is larger than the direct origin distances
    if (Math.abs(@distanceA - @distanceB) > @twoDistances.originDistances)
      console.error "Distances are impossibly far apart.", @distanceA, @distanceB
      throw new Error("Distances are impossibly far apart.", this)


  coordinates: ->
    if not @_coordinates?
      bottomLeftAngle = triangle.angleFromSides(
        @distanceB, @distanceA, @twoDistances.originDistances)
      y = triangle.side(bottomLeftAngle, @distanceA, 0.5 * Math.PI)
      x = (
        @twoDistances.originA.coordinates().x +
        rightTriangle.side(@distanceA, y))

      @_coordinates = new Coordinates(x, y)

    @_coordinates

exports.Graph =
class Graph
  constructor: (@originDistances, nodes) ->
    if not (@originDistances >= 0.0 && @originDistances < 1.0)
      throw new Error("originDistances=" + @originDistances)

    @originA = new Node(0.0, originDistances, null, 'rgba(255, 0, 0, 0.75)');
    @originA._coordinates = new Coordinates(0.0, 0);

    @originB = new Node(originDistances, 0.0, null, 'rgba(0, 0, 255, 0.75)')
    @originB._coordinates = new Coordinates(originDistances, 0);

    @nodes = [@originA, @originB];

    for node in nodes
      @nodes.push node

    for node in @nodes
      node.setTwoDistances @

    # The interior space filled by the graphic.
    @_xOffset = 0.5 * (1.0 - @originDistances)
    @_yOffset = 0

    @size = @_defaultSize
    @padding = @_defaultPadding

    @scale = @size - 2 * @padding

    @canvas = document.createElement('canvas')
    @canvas.width = @canvas.height = @size

    @graphics = @canvas.getContext('2d')

    @draw()

  _defaultSize: 256
  _defaultPadding: 16

  # Transform internal coordinates (relative to .originA)
  # into a canvas coordinates.
  transformX: (x) ->
    @padding + (@_xOffset + x) * @scale

  transformY: (y) ->
    @size - @padding - (y + @_yOffset) * @scale

  draw: ->
    @graphics.clearRect 0, 0, @width, @height

    for pass in [1, 2]
      for node in @nodes
        coordinates = node.coordinates();

        @graphics.strokeStyle = node.color || 'rgba(128, 128, 128, 1.0)';
        @graphics.fillStyle = node.color || 'rgba(0, 0, 0, 1.0)';
        @graphics.lineWidth = @scale / 128;

        # Drag the nodes after drawing the lines and borders.
        if pass == 2
            @graphics.beginPath()
            @graphics.arc(
                @transformX(coordinates.x),
                @transformY(coordinates.y),
                @scale / 32,
                0.0,
                2.0 * Math.PI,
                true)

            @graphics.fill()
            @graphics.stroke()
            @graphics.closePath()
            continue

        if node.source
            sourceCoordinates = node.source.coordinates()

            @graphics.beginPath()
            @graphics.moveTo(
                @transformX(coordinates.x),
                @transformY(coordinates.y))
            @graphics.lineTo(
                @transformX(sourceCoordinates.x),
                @transformY(sourceCoordinates.y))
            @graphics.stroke()
            @graphics.closePath()

        # XXX(JB): for testing, draw lines to origins
        if (node != @originA && node != @originB)
            @graphics.beginPath()
            @graphics.lineWidth = 0.25
            @graphics.strokeStyle = @originA.color
            @graphics.moveTo(
                @transformX(coordinates.x),
                @transformY(coordinates.y))
            @graphics.lineTo(
                @transformX(@originA.coordinates().x),
                @transformY(@originA.coordinates().y))
            @graphics.stroke()
            @graphics.closePath()

            @graphics.beginPath()
            @graphics.lineWidth = 0.25
            @graphics.strokeStyle = @originB.color
            @graphics.moveTo(
                @transformX(coordinates.x),
                @transformY(coordinates.y))
            @graphics.lineTo(
                @transformX(@originB.coordinates().x),
                @transformY(@originB.coordinates().y))
            @graphics.stroke()
            @graphics.closePath()

    @

# Basic trig utilities
triangle =
  angleFromSides: (opposite, adjacent1, adjacent2) ->
      # cosine law
      return Math.acos(
          (Math.pow(adjacent1, 2) + Math.pow(adjacent2, 2) - Math.pow(opposite, 2)
          ) / (2 * adjacent1 * adjacent2));

  side: (oppositeAngle, other, otherOppositeAngle) ->
      # sine law
      return (other / Math.sin(otherOppositeAngle)) * Math.sin(oppositeAngle);

rightTriangle =
  hypotenuse: (a, b) ->
    Math.sqrt(Math.pow(a, 2) + Math.pow(b, 2));

  side: (hypotenuse, opposite) ->
    Math.sqrt(Math.pow(hypotenuse, 2) - Math.pow(opposite, 2));
